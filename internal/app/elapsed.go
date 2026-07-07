package app

import (
	"context"
	"fmt"
	"time"

	"umamusume-notifier/internal/points"
)

func (m *Manager) SetElapsed(
	ctx context.Context,
	systemID string,
	minutes int,
) error {
	m.mu.Lock()

	system, reminder, ok := m.system(systemID)
	if !ok {
		m.mu.Unlock()
		return fmt.Errorf("unknown point system %q", systemID)
	}

	system.SetElapsed(time.Duration(minutes) * time.Minute)

	reminder.AlertSent = false
	reminder.FullSent = false

	systemToSave := system
	reminderToSave := reminder

	m.mu.Unlock()

	if err := m.store.SavePointSystems(
		ctx,
		[]*points.PointSystem{systemToSave},
	); err != nil {
		return fmt.Errorf("save point system: %w", err)
	}

	if err := m.store.SaveReminderState(
		ctx,
		reminderToSave,
	); err != nil {
		return fmt.Errorf("save reminder state: %w", err)
	}

	return nil
}