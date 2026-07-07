package app

import (
	"context"
	"time"

	"umamusume-notifier/internal/notification"
	"umamusume-notifier/internal/scheduler"
)

func (m *Manager) Tick(
	ctx context.Context,
	now time.Time,
) ([]notification.Event, error) {

	m.mu.Lock()
	defer m.mu.Unlock()

	events := m.tick(now)

	if err := m.persist(ctx); err != nil {
		return nil, err
	}

	return events, nil
}

func (m *Manager) tick(now time.Time) []notification.Event {
	var events []notification.Event

	for id, system := range m.pointSystems {
		reminder, ok := m.reminders[id]
		if !ok {
			panic("missing reminder state for point system " + id)
		}

		delta := now.Sub(system.LastTick)

		system.Advance(delta)
		system.LastTick = now

		event, ok := scheduler.ShouldAlert(
			system,
			reminder,
			m.alertThreshold,
		)
		if !ok {
			continue
		}

		switch event.Type {
		case notification.NearFull:
			reminder.AlertSent = true

		case notification.Full:
			reminder.AlertSent = true
			reminder.FullSent = true
		}

		events = append(events, event)
	}

	return events
}
