package app

import (
	"context"
	"fmt"

	"umamusume-notifier/internal/points"
)

func (m *Manager) persist(ctx context.Context) error {
	systems := make([]*points.PointSystem, 0, len(m.pointSystems))

	for id, system := range m.pointSystems {
		systems = append(systems, system)

		reminder := m.reminders[id]
		if reminder == nil {
			return fmt.Errorf("missing reminder state for point system %q", id)
		}

		if err := m.store.SaveReminderState(ctx, reminder); err != nil {
			return fmt.Errorf("save reminder state %q: %w", id, err)
		}
	}

	if err := m.store.SavePointSystems(ctx, systems); err != nil {
		return fmt.Errorf("save point systems: %w", err)
	}

	return nil
}
