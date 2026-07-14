# Umamusume Notifier Service — Development Guideline

A Telegram bot that tracks regenerating point systems (CP, TP, RP, ...), alerts when a system is close to full, and lets the user reply to consume points or correct the elapsed timer manually.

This document describes how to build the real Go service, based on the validated prototype behavior.

---

## 1. Project Structure

```
umamusume-notifier/
├── cmd/
│   └── server/
│       └── main.go              # entrypoint, wiring, graceful shutdown
├── internal/
│   ├── config/
│   │   └── config.go             # load + validate config.yaml
│   ├── points/
│   │   ├── model.go              # PointSystem struct, pure math
│   │   ├── manager.go            # in-memory state + mutex, orchestrates storage
│   │   └── scheduler.go          # ticker loop, threshold detection
│   ├── telegram/
│   │   ├── bot.go                # bot init, long-poll or webhook receiver
│   │   ├── commands.go           # /status /use /elapsed /help handlers
│   │   └── replies.go            # reply-to-message resolution
│   ├── storage/
│   │   ├── sqlite.go             # persistence
│   │   └── migrations/           # schema versioning
│   └── notifier/
│       └── reminder.go           # alert message composition + send
├── config.yaml
├── data.db
├── go.mod
└── Makefile
```

Keep `points` free of any Telegram or storage import — it should be testable as pure logic. `telegram` and `storage` depend on `points`, never the reverse.

---

## 2. Data Model

Do not store a continuously-incrementing float. Store **elapsed time toward the next point**, capped at `RegenMinutes`, matching the prototype:

```go
type PointSystem struct {
    ID           string        // "TP"
    Name         string        // "Training Points"
    Max          int
    Current      int
    RegenMinutes int           // minutes required to regenerate 1 point
    Elapsed      time.Duration // time accumulated toward the next point, capped at RegenMinutes
    LastTick     time.Time     // wall-clock time this system was last advanced
}
```

Reminder/alert state is tracked separately so it survives independently of point state:

```go
type ReminderState struct {
    SystemID       string
    AlertSent      bool
    FullSent       bool
    LastMessageID  int // Telegram message ID of the most recent reminder, for reply matching
}
```

### Why elapsed-based, not derived-from-clock-diff

The prototype settled on this approach for two reasons:

* It matches how a person mentally tracks progress ("TP is 4 minutes into a 10-minute regen"), which is what the `/elapsed` command manipulates directly.
* It composes cleanly with manual correction: setting `Elapsed` and re-running `Advance()` handles overflow (elapsed exceeding one regen interval awards multiple points) with no special-casing.

---

## 3. Regen Logic (`points/model.go`)

Single pure function, no I/O, easy to unit test:

```go
// Advance moves the system forward by delta and rolls elapsed time into
// points whenever it reaches RegenMinutes. Safe to call with any delta,
// including deltas larger than one full regen interval.
func (s *PointSystem) Advance(delta time.Duration) {
    if s.Current >= s.Max {
        s.Elapsed = 0
        return
    }
    s.Elapsed += delta
    regenDur := time.Duration(s.RegenMinutes) * time.Minute
    for s.Elapsed >= regenDur && s.Current < s.Max {
        s.Elapsed -= regenDur
        s.Current++
    }
    if s.Current >= s.Max {
        s.Current = s.Max
        s.Elapsed = 0
    }
}

// TimeUntilFull returns the duration remaining before Current reaches Max.
func (s *PointSystem) TimeUntilFull() time.Duration {
    if s.Current >= s.Max {
        return 0
    }
    regenDur := time.Duration(s.RegenMinutes) * time.Minute
    remaining := regenDur*time.Duration(s.Max-s.Current) - s.Elapsed
    if remaining < 0 {
        remaining = 0
    }
    return remaining
}
```

### Recovering from downtime

On startup, load `LastTick` from storage and call:

```go
sys.Advance(time.Since(sys.LastTick))
sys.LastTick = time.Now()
```

This replays any missed regeneration in one step — no need to loop minute-by-minute over the outage window.

### Consuming points

```go
func (s *PointSystem) Consume(amount int) {
    s.Current -= amount
    if s.Current < 0 {
        s.Current = 0
    }
    // Consuming keeps the in-progress regen timer intact.
}
```

### Manual elapsed correction (`/elapsed`)

```go
func (s *PointSystem) SetElapsed(minutes int) {
    s.Elapsed = 0
    s.Advance(time.Duration(minutes) * time.Minute) // reuses overflow handling
}
```

---

## 4. Scheduler (`points/scheduler.go`)

Runs on a `time.Ticker`, default interval 1 minute (configurable, independent of Telegram).

```go
for range ticker.C {
    now := time.Now()
    for _, sys := range manager.All() {
        delta := now.Sub(sys.LastTick)
        sys.Advance(delta)
        sys.LastTick = now
        checkThreshold(sys, reminderState[sys.ID])
    }
    manager.PersistAll() // batch write, not per-system
}
```

`checkThreshold` mirrors the prototype exactly:

* `Current >= Max` and `!FullSent` → send full notification, mark `FullSent = true`, `AlertSent = true` (prevents a stale near-full alert firing after full).
* `TimeUntilFull() <= AlertThreshold` and `!AlertSent` → send near-full alert, mark `AlertSent = true`.
* Both flags reset to `false` whenever `Consume` or `SetElapsed` is called on that system — this is what lets the bot alert again on the next cycle.

Keep the threshold check and the send call decoupled (`shouldAlert(sys) bool` as a pure function) so it's testable without a real Telegram client.

---

## 5. Telegram Integration

### Library

Use [`go-telegram-bot-api/telegram-bot-api`](https://github.com/go-telegram-bot-api/telegram-bot-api) (or `gotgbot` if you want typed handlers) with **long polling** for a single-user personal bot — no need for webhook/HTTPS infrastructure. Switch to webhook mode only if you deploy behind a domain with TLS.

### Reply-to-message resolution

This is the core UX mechanic and needs to be exact:

1. When a reminder is sent, store the returned `message_id` in `ReminderState.LastMessageID` for that system.
2. Every incoming Telegram update, check `update.Message.ReplyToMessage`. If present, look up which system owns that `message_id`.
3. If the reply body parses as a number, treat it as `Consume(amount)` for that system.
4. If no `ReplyToMessage` is present, fall back to slash commands only — do not guess which system a bare number refers to.

```go
func resolveReplyTarget(state map[string]*ReminderState, replyToID int) (systemID string, ok bool) {
    for id, s := range state {
        if s.LastMessageID == replyToID {
            return id, true
        }
    }
    return "", false
}
```

### Commands

| Command | Behavior |
|---|---|
| `/status` | List all systems: current/max, time to full or FULL |
| `/use SYSTEM amount` | Immediate consume, also accept negative amount to add point instead |
| `/elapsed SYSTEM minutes` | Manually set elapsed time toward next point (overflow rolls into points) |
| `/help` | List commands with short descriptions |
| plain number, replying to an alert | Consume that amount from the alert's system |

Validate `SYSTEM` against known IDs and reject unparseable amounts with a short usage message — do not silently ignore malformed input, since there's no UI feedback loop otherwise.

### Authorization

Since this is a personal bot, hard-restrict handling to a single `chat_id` from config. Ignore or log-and-drop updates from any other chat ID — do not process commands from arbitrary users even if they discover the bot.

---

## 6. Persistence (SQLite)

```sql
CREATE TABLE systems (
    id             TEXT PRIMARY KEY,
    name           TEXT NOT NULL,
    max_points     INTEGER NOT NULL,
    current_points INTEGER NOT NULL,
    regen_minutes  INTEGER NOT NULL,
    elapsed_seconds INTEGER NOT NULL DEFAULT 0,
    last_tick      DATETIME NOT NULL
);

CREATE TABLE reminder_state (
    system_id        TEXT PRIMARY KEY REFERENCES systems(id),
    alert_sent       BOOLEAN NOT NULL DEFAULT 0,
    full_sent        BOOLEAN NOT NULL DEFAULT 0,
    last_message_id  INTEGER
);
```

* Write on every state-changing event (consume, elapsed-set) immediately — these are rare and user-initiated, so latency doesn't matter.
* Write from the scheduler tick in a batch, not per-system, to avoid excessive fsyncs on a 1-minute cadence.
* Use `modernc.org/sqlite` (pure Go, no cgo) unless you have a reason to need `mattn/go-sqlite3`.

---

## 7. Configuration

```yaml
telegram:
  token: "xxxxx"
  chat_id: 123456789

scheduler:
  tick_interval: 1m
  alert_threshold: 30m

systems:
  - id: CP
    name: "Combat Points"
    max: 1
    regen_minutes: 480
  - id: TP
    name: "Training Points"
    max: 100
    regen_minutes: 10
  - id: RP
    name: "Raid Points"
    max: 5
    regen_minutes: 120
```

`current_points` and `elapsed` are **not** in config — they're runtime state, seeded once on first boot (current=0, elapsed=0) and owned by SQLite thereafter. Adding a new point system is purely a config change plus a migration to insert its initial row.

---

## 8. Concurrency

* One `sync.RWMutex` (or a single-goroutine actor pattern) guards the in-memory `map[string]*PointSystem`. The scheduler tick and Telegram command handlers both mutate this map — never let them race.
* Prefer an actor/channel pattern over raw mutexes if the codebase grows past a handful of operations: a single goroutine owns state, all reads/writes go through a channel of commands. This avoids subtle lock-ordering bugs and keeps `points` easy to unit test in isolation.

---

## 9. Testing

* `points` package: table-driven tests for `Advance`, `Consume`, `SetElapsed`, `TimeUntilFull`, including overflow cases (elapsed far exceeding one regen interval, consuming more than current, elapsed set on an already-full system).
* `scheduler`: test `shouldAlert` as a pure function against fabricated `PointSystem` + `ReminderState` combinations — do not spin up a real ticker in unit tests.
* `telegram`: test reply-resolution and command parsing against fixture `Update` payloads; do not hit the real Bot API in tests. Use an interface (`type Sender interface { SendMessage(...) (int, error) }`) so a mock sender can be injected.

---

## 10. Suggested Build Order

1. `points` package + tests (pure logic, no dependencies) — this is the part the prototype already validated behaviorally.
2. `storage` package + migration, wire up load/persist for `points`.
3. `scheduler` running against in-memory state, logging alerts to stdout (no Telegram yet).
4. `telegram` bot: `/status` and `/help` first (read-only), then `/use`, `/elapsed`, then reply-based consume.
5. Deploy as a single long-running binary (systemd service); no need for orchestration at this scale.

---

## 11. Others

* Time zone handling — prototype uses a simulated clock; real service should store and compute in UTC internally, only converting to local time for display.
* Telegram API failures / rate limits — add retry with backoff on `SendMessage`, and don't crash the scheduler tick if one send fails.
* Config hot-reload — not required initially; restart the service to add a new point system.
* Multi-user support — out of scope per the single-`chat_id` authorization model above.
