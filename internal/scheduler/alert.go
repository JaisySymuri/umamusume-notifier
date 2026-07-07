package scheduler

import (
	"time"

	"umamusume-notifier/internal/notification"
	"umamusume-notifier/internal/points"
)

func ShouldAlert(
	system *points.PointSystem,
	reminder *points.ReminderState,
	threshold time.Duration,
) (notification.Event, bool) {

	if system.Current >= system.Max {
		if reminder.FullSent {
			return notification.Event{}, false
		}

		return notification.Event{
			SystemID: system.ID,
			Type:     notification.Full,
		}, true
	}

	if reminder.AlertSent {
		return notification.Event{}, false
	}

	if system.TimeUntilFull() <= threshold {
		return notification.Event{
			SystemID: system.ID,
			Type:     notification.NearFull,
		}, true
	}

	if system == nil || reminder == nil {
		return notification.Event{}, false
	}

	return notification.Event{}, false
}
