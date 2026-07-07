package storage

import (
	"context"
	"fmt"
	"time"

	"umamusume-notifier/internal/points"
)

// LoadPointSystems loads all persisted point systems.
func (s *SQLiteStore) LoadPointSystems(
	ctx context.Context,
) ([]*points.PointSystem, error) {

	rows, err := s.db.QueryContext(ctx, loadPointSystemsQuery)
	if err != nil {
		return nil, fmt.Errorf("query point systems: %w", err)
	}
	defer rows.Close()

	var systems []*points.PointSystem

	for rows.Next() {
		var (
			ps             points.PointSystem
			elapsedSeconds int64
			lastTick       time.Time
		)

		if err := rows.Scan(
			&ps.ID,
			&ps.Name,
			&ps.Max,
			&ps.Current,
			&ps.RegenMinutes,
			&elapsedSeconds,
			&lastTick,
		); err != nil {
			return nil, fmt.Errorf("scan point system: %w", err)
		}

		ps.Elapsed = time.Duration(elapsedSeconds) * time.Second
		ps.LastTick = lastTick.UTC()

		systems = append(systems, &ps)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate point systems: %w", err)
	}

	return systems, nil
}

// SavePointSystems persists runtime state for all point systems.
func (s *SQLiteStore) SavePointSystems(
	ctx context.Context,
	systems []*points.PointSystem,
) error {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.PrepareContext(ctx, savePointSystemQuery)
	if err != nil {
		return fmt.Errorf("prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, ps := range systems {
		if ps == nil {
			continue
		}
		if _, err := stmt.ExecContext(
			ctx,
			ps.Current,
			int64(ps.Elapsed/time.Second),
			ps.LastTick.UTC(),
			ps.ID,
		); err != nil {
			return fmt.Errorf("update point system %q: %w", ps.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
