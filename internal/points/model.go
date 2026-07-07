package points

import (
	"time"
)

// PointSystem represents the runtime state of a regenerating point system.
type PointSystem struct {
	Definition

	Current int

	// Progress accumulated toward the next regenerated point.
	Elapsed time.Duration

	// Last time this system was updated.
	LastTick time.Time
}

// ReminderState stores notification state separately from PointSystem.
type ReminderState struct {
	SystemID string

	AlertSent bool
	FullSent  bool

	// Telegram message ID of the latest reminder.
	LastMessageID int
}


// New creates a validated PointSystem from a Definition.
func New(def Definition) (*PointSystem, error) {
	switch {
	case def.ID == "":
		return nil, ErrEmptyID

	case def.Name == "":
		return nil, ErrEmptyName

	case def.Max <= 0:
		return nil, ErrInvalidMax

	case def.RegenMinutes <= 0:
		return nil, ErrInvalidRegenTime
	}

	return &PointSystem{
		Definition: def,

		Current:  0,
		Elapsed:  0,
		LastTick: time.Now().UTC(),
	}, nil
}