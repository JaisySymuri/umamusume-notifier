# Architecture Decisions

This document records architectural decisions made during implementation that intentionally differ from the original development guideline.

The guideline remains the overall roadmap, while this document explains refinements made during implementation.

---

# ADR-001: Keep `points` as a Pure Domain Package

## Original Guideline

The guideline places:

```text
points/
    model.go
    manager.go
    scheduler.go
```

inside the same package.

## Decision

`points` contains only domain logic.

It does not know about:

* SQLite
* Telegram
* scheduler orchestration
* mutexes
* configuration
* logging

The package only models regenerating point systems.

## Reason

This keeps the business logic:

* deterministic
* easy to unit test
* reusable
* independent from infrastructure

---

# ADR-002: Introduce an `app` Layer

Instead of allowing `points` to coordinate storage, persistence, and scheduler behavior, orchestration lives in a separate `app` package.

```text
internal/
    app/
    notification/
    points/
    scheduler/
    storage/
```

## Reason

The manager is an application concern, not a domain concern.

`app.Manager` owns in-memory state, locking, persistence calls, and alert-state transitions.

---

# ADR-003: Split Configuration from Runtime State

The original guideline combines immutable configuration and mutable runtime state inside `PointSystem`.

Implementation introduces:

```text
Definition
```

for immutable configuration and

```text
PointSystem
```

for runtime state.

Conceptually:

```text
PointSystem = Definition + Runtime State
```

`PointSystem` embeds `Definition`.

## Reason

Configuration and runtime evolve independently.

Configuration:

* ID
* Name
* Max
* RegenMinutes

Runtime:

* Current
* Elapsed
* LastTick

This separation makes storage synchronization much cleaner.

---

# ADR-004: Use a Shared Notification Contract

The original guideline does not define a notification-event package.

Implementation introduces a small `internal/notification` package:

```text
notification.Type
notification.Event
```

## Reason

`app` and `scheduler` both need to talk about alert outcomes, but neither package should own the cross-package event contract.

Keeping the event type in `internal/notification` avoids circular imports and keeps the scheduler focused on detection while `app` focuses on state mutation and persistence.

---

# ADR-005: Synchronize Configuration Explicitly

The original guideline does not explicitly define how configured systems become database rows.

Implementation introduces:

```text
SyncPointSystems()
```

Application startup becomes:

```text
Load config
-> Initialize database
-> SyncPointSystems()
-> LoadPointSystems()
-> Start scheduler
```

## Reason

This separates two concerns:

* Configuration defines what systems exist.
* SQLite stores mutable runtime state.

Batch persistence is used for scheduler-driven updates, while user-triggered state changes can still be persisted immediately.

---

# ADR-006: Store Interface

The storage package exposes a small interface:

```text
Store
```

instead of directly exposing SQLite.

The application depends on:

```text
storage.Store
```

not

```text
SQLiteStore
```

## Reason

This keeps higher layers independent from the storage implementation.

Alternative implementations (PostgreSQL, JSON, memory) can be introduced later.

The application persists all point systems in one batch call rather than saving each system individually during a tick.

---

# ADR-007: Storage Owns Initialization

New point systems are initialized by the storage layer.

Initial runtime values:

* Current = 0
* Elapsed = 0
* LastTick = current UTC time

Configuration only describes the point system.

It does not specify runtime state.

---

# ADR-008: Runtime Persistence Only

Runtime updates only modify:

* Current
* Elapsed
* LastTick
* reminder flags such as `AlertSent` and `FullSent`

Configuration fields are synchronized only through `SyncPointSystems()`.

This prevents scheduler updates from accidentally modifying immutable configuration.

Reminder state is persisted alongside runtime state, but it remains separate from the point-system row so alert flags can be reset independently.

---

# ADR-009: Use `time.Duration` Inside the Domain

The original guideline accepts elapsed time as minutes.

Implementation uses:

```text
time.Duration
```

throughout the domain package.

Conversion from minutes happens only at the interface layer (Telegram commands, configuration, and persistence boundaries).

## Reason

The domain should operate on native Go types rather than user-facing units.

---

# ADR-010: Constructor Uses a Definition

Instead of:

```text
New(
    id,
    name,
    max,
    regenMinutes,
)
```

implementation uses:

```text
New(Definition)
```

## Reason

This avoids positional parameter mistakes and naturally grows if additional configuration fields are added later.

---

# ADR-011: SQL Is Isolated

All SQL statements are defined in `schema.go`.

Implementation files contain only Go logic.

## Reason

Separating SQL from business code improves readability and makes future migrations easier to manage.

---

# ADR-012: Testing Strategy

The project follows two layers of testing.

## Unit Tests

The `points` package is tested in isolation with table-driven tests.

Goal:

* verify business rules
* no external dependencies

## Integration Tests

The `storage` package uses a temporary SQLite database to verify:

* schema creation
* synchronization
* persistence
* loading
* reminder state

This verifies the complete storage layer without requiring Telegram or the scheduler.

The scheduler/application boundary is also kept narrow enough to unit test by mocking the `TickManager` interface instead of depending on a concrete app manager.

---

# ADR-013: Keep `notification` Minimal

The `notification` package only defines alert outcomes shared between `app` and `scheduler`.

It does not contain:

* transport logic
* message formatting
* Telegram API calls
* persistence code

## Reason

This keeps notification semantics simple and reusable while leaving delivery concerns to the eventual Telegram or bot layer.

---

# ADR-014: `app` Owns Orchestration

`app.Manager` coordinates the runtime map of point systems, reminder state, persistence, and scheduler-triggered advancement.

## Reason

The application layer is the right place for:

* locking
* multi-system iteration
* loading and saving runtime state
* applying notification side effects after threshold checks

That logic is broader than the pure `points` domain, but narrower than a full transport layer.


Current structure: Before phase 4 started


```
📦umamusume-notifier
 ┣ 📂internal
 ┃ ┣ 📂app
 ┃ ┃ ┣ 📜manager.go
 ┃ ┃ ┣ 📜persist.go
 ┃ ┃ ┣ 📜tick.go
 ┃ ┃ ┗ 📜types.go
 ┃ ┣ 📂notification
 ┃ ┃ ┗ 📜types.go
 ┃ ┣ 📂points
 ┃ ┃ ┣ 📜advance.go
 ┃ ┃ ┣ 📜advance_test.go
 ┃ ┃ ┣ 📜consume.go
 ┃ ┃ ┣ 📜consume_test.go
 ┃ ┃ ┣ 📜definition.go
 ┃ ┃ ┣ 📜doc.go
 ┃ ┃ ┣ 📜elapsed.go
 ┃ ┃ ┣ 📜elapsed_test.go
 ┃ ┃ ┣ 📜errors.go
 ┃ ┃ ┣ 📜model.go
 ┃ ┃ ┣ 📜model_test.go
 ┃ ┃ ┣ 📜test_helpers_test.go
 ┃ ┃ ┣ 📜time.go
 ┃ ┃ ┗ 📜time_test.go
 ┃ ┣ 📂scheduler
 ┃ ┃ ┣ 📜alert.go
 ┃ ┃ ┗ 📜scheduler.go
 ┃ ┗ 📂storage
 ┃ ┃ ┣ 📜doc.go
 ┃ ┃ ┣ 📜migrate.go
 ┃ ┃ ┣ 📜point_system.go
 ┃ ┃ ┣ 📜reminder_state.go
 ┃ ┃ ┣ 📜repository.go
 ┃ ┃ ┣ 📜schema.go
 ┃ ┃ ┣ 📜sqlite.go
 ┃ ┃ ┣ 📜storage_integration_test.go
 ┃ ┃ ┗ 📜sync.go
 ┣ 📜architecture-decisions.md
 ┣ 📜coverage
 ┣ 📜go.mod
 ┣ 📜go.sum
 ┗ 📜umamusume-notifier-dev-guideline.md

 
```


## Proposed Phase 4 roadmap

1. **Create `internal/telegram` package.**
2. Wire up the bot with long polling and single-`chat_id` authorization.
3. Implement `/help` (static output).
4. Implement `/status` (read-only via `app.Manager.Status()`).
5. Add `formatter.go` and unit-test message formatting.
6. Implement `/use`.
7. Implement `/elapsed`.
8. Implement reply-based consume using stored `LastMessageID`.
9. Finally, connect scheduler notification events to the Telegram sender so alerts are delivered automatically.
