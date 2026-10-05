package app

import (
	"context"
	"time"

	"umamusume-notifier/internal/metrics"
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

		previousPoints := system.Current
		delta := now.Sub(system.LastTick)

		system.Advance(delta)
		system.LastTick = now

		if id == "LP" {
			firstMilestone := (previousPoints/3 + 1) * 3
			if previousPoints < 0 {
				firstMilestone = 0
			}
			for points := firstMilestone; points <= system.Current && points > previousPoints; points += 3 {
				events = append(events, notification.Event{
					SystemID:    system.ID,
					SystemName:  system.Name,
					Type:        notification.Milestone,
					Points:      points,
					ScheduledAt: now,
				})
				metrics.ObserveReminderScheduled()
			}
		}

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
			if reminder.FullSince.IsZero() {
				reminder.FullSince = now
			}
		}

		event.ScheduledAt = now
		metrics.ObserveReminderScheduled()
		events = append(events, event)
	}

	return events
}
