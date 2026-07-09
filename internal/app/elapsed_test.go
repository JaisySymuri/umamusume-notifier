package app

import (
	"context"
	"testing"
	"time"

	"umamusume-notifier/internal/points"
)

func TestManagerSetElapsed(t *testing.T) {
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
			},
		},
		reminders: map[string]*points.ReminderState{
			"TP": {
				SystemID:  "TP",
				AlertSent: true,
				FullSent:  true,
			},
		},
	}

	if err := manager.SetElapsed(
		context.Background(),
		"TP",
		5,
	); err != nil {
		t.Fatalf("SetElapsed() error = %v", err)
	}

	system := manager.pointSystems["TP"]

	if system.Current != 50 {
		t.Fatalf("Current = %d, want 50", system.Current)
	}

	if system.Elapsed != 5*time.Minute {
		t.Fatalf(
			"Elapsed = %v, want %v",
			system.Elapsed,
			5*time.Minute,
		)
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

func TestManagerSetElapsed_UnknownSystem(t *testing.T) {
	manager := &Manager{
		store: &mockStore{},
	}

	if err := manager.SetElapsed(
		context.Background(),
		"UNKNOWN",
		10,
	); err == nil {
		t.Fatal("expected error")
	}
}

func TestManagerSetRegen(t *testing.T) {
	store := &mockStore{}

	manager := &Manager{
		store: store,
		pointSystems: map[string]*points.PointSystem{
			"TP": {
				Definition: points.Definition{
					ID:           "TP",
					Name:         "Training Points",
					Max:          100,
					RegenMinutes: 10,
				},
				Current: 50,
				Elapsed: 4 * time.Minute,
			},
		},
		reminders: map[string]*points.ReminderState{
			"TP": {
				SystemID:  "TP",
				AlertSent: true,
				FullSent:  true,
			},
		},
	}

	if err := manager.SetRegen(
		context.Background(),
		"TP",
		6,
	); err != nil {
		t.Fatalf("SetRegen() error = %v", err)
	}

	system := manager.pointSystems["TP"]

	if system.Elapsed != 4*time.Minute {
		t.Fatalf("Elapsed = %v, want %v", system.Elapsed, 4*time.Minute)
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
