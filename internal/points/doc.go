// Package points contains the core domain logic for regenerating point systems.
//
// The package is intentionally independent from storage, scheduling, and
// Telegram so it can be tested as pure business logic.
//
// All time calculations are based on elapsed durations rather than periodic
// ticking, allowing the system to recover from downtime by replaying the
// elapsed time since the last update.
package points