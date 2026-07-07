package app

import (
	"context"
	"fmt"
)

func (m *Manager) RecordReminderMessage(
	ctx context.Context,
	systemID string,
	messageID int,
) error {
	m.mu.Lock()

	_, reminder, ok := m.system(systemID)
	if !ok {
		m.mu.Unlock()
		return fmt.Errorf("unknown point system %q", systemID)
	}

	reminder.LastMessageID = messageID

	reminderToSave := reminder

	m.mu.Unlock()

	if err := m.store.SaveReminderState(
		ctx,
		reminderToSave,
	); err != nil {
		return fmt.Errorf("save reminder state: %w", err)
	}

	return nil
}