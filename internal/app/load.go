package app

import (
	"context"
	"fmt"

	"umamusume-notifier/internal/points"
)

// Load synchronizes configured point systems with storage and loads
// persisted runtime state into memory.
func (m *Manager) Load(
	ctx context.Context,
	definitions []points.Definition,
) error {
	if err := m.store.SyncPointSystems(ctx, definitions); err != nil {
		return fmt.Errorf("sync point systems: %w", err)
	}

	systems, err := m.store.LoadPointSystems(ctx)
	if err != nil {
		return fmt.Errorf("load point systems: %w", err)
	}

	reminders, err := m.store.LoadReminderStates(ctx)
	if err != nil {
		return fmt.Errorf("load reminder states: %w", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.pointSystems = make(map[string]*points.PointSystem, len(systems))
	for _, system := range systems {
		m.pointSystems[system.ID] = system
	}

	m.reminders = make(map[string]*points.ReminderState, len(reminders))
	for _, reminder := range reminders {
		m.reminders[reminder.SystemID] = reminder
	}

	for id := range m.pointSystems {
		if _, ok := m.reminders[id]; ok {
			continue
		}

		m.reminders[id] = &points.ReminderState{
			SystemID: id,
		}
	}

	return nil
}
