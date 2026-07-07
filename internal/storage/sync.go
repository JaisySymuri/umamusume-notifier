package storage

import (
	"context"
	"fmt"
	"time"

	"umamusume-notifier/internal/points"
)

type Definition struct {
	ID           string
	Name         string
	Max          int
	RegenMinutes int
}

type PointSystem struct {
	Definition

	Current  int
	Elapsed  time.Duration
	LastTick time.Time
}

const syncPointSystemQuery = `
INSERT INTO point_systems (
    id,
    name,
    max_points,
    current_points,
    regen_minutes,
    elapsed_seconds,
    last_tick
)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    name = excluded.name,
    max_points = excluded.max_points,
    regen_minutes = excluded.regen_minutes;
`

// SyncPointSystems synchronizes configured point systems into storage.
func (s *SQLiteStore) SyncPointSystems(
	ctx context.Context,
	definitions []points.Definition,
) error {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.PrepareContext(ctx, syncPointSystemQuery)
	if err != nil {
		return fmt.Errorf("prepare sync statement: %w", err)
	}
	defer stmt.Close()

	now := time.Now().UTC()

	for _, def := range definitions {
		_, err := stmt.ExecContext(
			ctx,
			def.ID,
			def.Name,
			def.Max,

			0, // initial current
			def.RegenMinutes,
			0,   // initial elapsed
			now, // initial last tick
		)
		if err != nil {
			return fmt.Errorf("sync point system %q: %w", def.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
