package app

import (
	"time"

	"umamusume-notifier/internal/points"
)

func resetFullTracking(reminder *points.ReminderState) {
	reminder.AlertSent = false
	reminder.FullSent = false
	reminder.FullSince = time.Time{}
	reminder.FullOverHourSent = false
}

func (m *Manager) logManualAdjustment(system *points.PointSystem, reminder *points.ReminderState, action string, now time.Time) {
	if m.logger == nil || reminder == nil || reminder.FullSince.IsZero() {
		return
	}

	fullFor := now.Sub(reminder.FullSince)
	if fullFor < 0 {
		fullFor = 0
	}

	fullForMinutes := int(fullFor.Minutes())
	lateMinutes := fullForMinutes - 60
	if lateMinutes < 0 {
		lateMinutes = 0
	}
	late := lateMinutes > 0

	m.logger.Printf(
		"event=manual_adjust_after_full action=%s system_id=%s system_name=%q full_for_minutes=%d late=%t late_minutes=%d adjusted_at=%s",
		action,
		system.ID,
		system.Name,
		fullForMinutes,
		late,
		lateMinutes,
		now.UTC().Format(time.RFC3339),
	)
}
