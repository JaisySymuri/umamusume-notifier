package storage

import (
	"context"

	"umamusume-notifier/internal/points"
)

// Store defines the persistence layer for the application.
//
// Implementations are responsible for loading and saving runtime state.
// Configuration (such as point definitions) is owned by the config package.
type Store interface {
	// Initialize prepares the underlying storage for use.
	Initialize(ctx context.Context) error

	// SyncPointSystems synchronizes configured point systems into storage.
	// Missing systems are inserted while immutable metadata is kept up-to-date.
	// Runtime state is preserved for existing systems.
	SyncPointSystems(
        ctx context.Context,
        definitions []points.Definition,
    ) error

	// LoadPointSystems loads persisted point systems.
	LoadPointSystems(ctx context.Context) ([]*points.PointSystem, error)

	// SavePointSystems persists runtime state for all point systems.
	SavePointSystems(ctx context.Context, systems []*points.PointSystem) error

	// LoadReminderStates loads reminder state for every point system.
	LoadReminderStates(ctx context.Context) ([]*points.ReminderState, error)

	// SaveReminderState persists a single reminder state.
	SaveReminderState(ctx context.Context, state *points.ReminderState) error
}