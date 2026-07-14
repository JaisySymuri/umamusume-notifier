package app

import (
	"context"
	"testing"
	"time"

	"umamusume-notifier/internal/points"
)

func TestManagerConsumeReply(t *testing.T) {
	store := &mockStore{}

	manager := &Manager{
		store: store,
		pointSystems: map[string]*points.PointSystem{
			"TP": {
				Definition: points.Definition{
					ID:            "TP",
					Name:          "Training Points",
					Max:           100,
					RegenMinutes:  10,
				},
				Current: 50,
				Elapsed: 5 * time.Minute,
			},
		},
		reminders: map[string]*points.ReminderState{
			"TP": {
				SystemID:      "TP",
				AlertSent:     true,
				FullSent:      true,
				LastMessageID: 123,
			},
		},
	}

	err := manager.ConsumeReply(
		context.Background(),
		123,
		10,
	)
	if err != nil {
		t.Fatalf("ConsumeReply() error = %v", err)
	}

	system := manager.pointSystems["TP"]

	if system.Current != 40 {
		t.Fatalf("Current = %d, want 40", system.Current)
	}

	if system.Elapsed != 5*time.Minute {
		t.Fatalf("Elapsed = %v, want 5m0s", system.Elapsed)
	}

	reminder := manager.reminders["TP"]

	if reminder.AlertSent {
		t.Fatal("AlertSent should be reset")
	}

	if reminder.FullSent {
		t.Fatal("FullSent should be reset")
	}

	if !store.savePointSystemsCalled {
		t.Fatal("SavePointSystems was not called")
	}

	if !store.saveReminderCalled {
		t.Fatal("SaveReminderState was not called")
	}
}

func TestManagerConsumeReply_UnknownMessage(t *testing.T) {
	manager := &Manager{
		store: &mockStore{},
		reminders: map[string]*points.ReminderState{
			"TP": {
				SystemID:      "TP",
				LastMessageID: 999,
			},
		},
	}

	err := manager.ConsumeReply(
		context.Background(),
		123,
		10,
	)

	if err == nil {
		t.Fatal("expected error")
	}
}
