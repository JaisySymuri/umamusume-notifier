package storage

import (
	"context"
	"fmt"
)

// Initialize prepares the SQLite database for use.
func (s *SQLiteStore) Initialize(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		// Rollback is safe even if Commit has already succeeded.
		_ = tx.Rollback()
	}()

	if _, err := tx.ExecContext(ctx, createPointSystemsTable); err != nil {
		return fmt.Errorf("create point_systems table: %w", err)
	}

	if _, err := tx.ExecContext(ctx, createReminderStatesTable); err != nil {
		return fmt.Errorf("create reminder_states table: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}