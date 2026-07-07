package storage

import (
	"context"
	"database/sql"
	"fmt"

	"umamusume-notifier/internal/points"
)

// LoadReminderStates loads all reminder states.
func (s *SQLiteStore) LoadReminderStates(
	ctx context.Context,
) ([]*points.ReminderState, error) {

	rows, err := s.db.QueryContext(ctx, loadReminderStatesQuery)
	if err != nil {
		return nil, fmt.Errorf("query reminder states: %w", err)
	}
	defer rows.Close()

	var states []*points.ReminderState

	for rows.Next() {
		var (
			state         points.ReminderState
			lastMessageID sql.NullInt64
		)

		if err := rows.Scan(
			&state.SystemID,
			&state.AlertSent,
			&state.FullSent,
			&lastMessageID,
		); err != nil {
			return nil, fmt.Errorf("scan reminder state: %w", err)
		}

		if lastMessageID.Valid {
			state.LastMessageID = int(lastMessageID.Int64)
		}

		states = append(states, &state)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate reminder states: %w", err)
	}

	return states, nil
}

// SaveReminderState persists a reminder state.
func (s *SQLiteStore) SaveReminderState(
	ctx context.Context,
	state *points.ReminderState,
) error {

	if state == nil {
		return nil
	}

	_, err := s.db.ExecContext(
		ctx,
		saveReminderStateQuery,
		state.SystemID,
		state.AlertSent,
		state.FullSent,
		state.LastMessageID,
	)
	if err != nil {
		return fmt.Errorf("save reminder state %q: %w", state.SystemID, err)
	}

	return nil
}