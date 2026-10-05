package telegram

import (
	"strings"
	"testing"

	"umamusume-notifier/internal/notification"
)

func TestFormatNotificationLivePointMilestone(t *testing.T) {
	got := FormatNotification(notification.Event{
		SystemID: "LP",
		Type:     notification.Milestone,
		Points:   6,
	})
	if !strings.Contains(got, "Live Points") || !strings.Contains(got, "6") {
		t.Fatalf("FormatNotification() = %q, want Live Points milestone message", got)
	}
}
