# Umamusume Notifier

Umamusume Notifier is a Telegram bot for tracking regenerating point systems like CP, TP, and RP.

It keeps the current points in SQLite, sends reminders when a system is near full or full, and lets you update the state directly from Telegram commands or by replying to a reminder message.

## Features

- Track multiple point systems at once
- Send Telegram reminders when a system is near full or full
- Reply to a reminder message with a number to consume points
- Use slash commands to check status, consume points, set points, and correct timer state
- Persist state in SQLite

## Project Structure

- `cmd/server` - application entrypoint
- `internal/app` - orchestration and persistence coordination
- `internal/config` - YAML config loading and validation
- `internal/points` - pure domain logic for point regeneration
- `internal/scheduler` - alert detection and ticking
- `internal/storage` - SQLite storage implementation
- `internal/telegram` - Telegram bot commands, replies, and formatting

## Requirements

- Go 1.22 or newer
- A Telegram bot token from BotFather
- A Telegram chat ID for the private chat where the bot should respond

## Configuration

Copy `config.example.yaml` to `config.yaml` and fill in your values.

Example:

```yaml
telegram:
  token: "123456:ABCDEF"
  chat_id: 123456789

scheduler:
  tick_interval: 1m
  alert_threshold: 30m

systems:
  - id: CP
    name: "Club Points"
    max: 1
    regen_minutes: 480
  - id: TP
    name: "Training Points"
    max: 100
    regen_minutes: 10
  - id: RP
    name: "Race Points"
    max: 5
    regen_minutes: 120
```

### Config Fields

- `telegram.token` - Telegram bot token
- `telegram.chat_id` - only this chat ID is allowed to use the bot
- `scheduler.tick_interval` - how often the scheduler advances point regeneration
- `scheduler.alert_threshold` - when to send the "almost full" reminder
- `systems` - list of point systems to track

## Running

1. Make sure `config.yaml` exists in the project root.
2. Run the service:

```bash
go run ./cmd/server
```

The app creates and uses `data.db` in the project root.

## Telegram Commands

- `/status` - show all point systems
- `/help` - show help text
- `/use <SYSTEM> <AMOUNT>` - consume points, or add points if `AMOUNT` is negative
- `/set <SYSTEM> <AMOUNT>` - set the current points directly
- `/elapsed <SYSTEM> <MINUTES>` - set elapsed regeneration time
- `/regen <SYSTEM> <MINUTES_LEFT>` - set how many minutes remain until the next point

### Reply Flow

When the bot sends a reminder like:

```text
✅ TP is now full.

Reply with the amount of points you use.
```

you can reply directly with a number such as `60`.

That reply is mapped back to the reminder message and processed as a point consumption for the matching system.

## Behavior Notes

- The bot is restricted to one chat ID from config.
- Point values are persisted in SQLite.
- Reminder flags are reset after a manual consume or point set.
- Manual consume keeps any in-progress regeneration time intact.
- Current points are clamped between `0` and the system maximum.

## Development

Run the tests:

```bash
go test ./...
```

## Design

The codebase separates responsibilities clearly:

- `points` contains pure domain logic
- `app` coordinates state, locking, and persistence
- `telegram` handles Telegram-specific input/output
- `storage` owns the database layer

This keeps the regeneration logic easy to test and the transport layers easy to change.
