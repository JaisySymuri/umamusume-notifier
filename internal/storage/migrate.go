package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"umamusume-notifier/internal/metrics"
)

// Initialize prepares the SQLite database for use.
func (s *SQLiteStore) Initialize(ctx context.Context) error {
	start := time.Now()
	defer metrics.ObserveStorageOp("initialize", time.Since(start))

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

	if err := ensureReminderStateColumns(ctx, tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func ensureReminderStateColumns(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(
		ctx,
		`ALTER TABLE reminder_states ADD COLUMN full_since DATETIME;`,
	); err != nil && !columnExistsError(err) {
		return fmt.Errorf("add reminder_states.full_since: %w", err)
	}

	if _, err := tx.ExecContext(
		ctx,
		`ALTER TABLE reminder_states ADD COLUMN full_over_hour_sent BOOLEAN NOT NULL DEFAULT FALSE;`,
	); err != nil && !columnExistsError(err) {
		return fmt.Errorf("add reminder_states.full_over_hour_sent: %w", err)
	}

	return nil
}

func columnExistsError(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "duplicate column name") || strings.Contains(err.Error(), "already exists"))
}
