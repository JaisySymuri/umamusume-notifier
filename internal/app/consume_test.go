package app

import (
	"context"
	"testing"
	"time"

	"umamusume-notifier/internal/points"
)

type mockStore struct {
	savePointSystemsCalled bool
	saveReminderCalled     bool
}

func (m *mockStore) Initialize(context.Context) error {
	return nil
}

func (m *mockStore) SyncPointSystems(context.Context, []points.Definition) error {
	return nil
}

func (m *mockStore) LoadPointSystems(context.Context) ([]*points.PointSystem, error) {
	return nil, nil
}

func (m *mockStore) LoadReminderStates(context.Context) ([]*points.ReminderState, error) {
	return nil, nil
}

func (m *mockStore) SavePointSystems(
	context.Context,
	[]*points.PointSystem,
) error {
	m.savePointSystemsCalled = true
	return nil
}

func (m *mockStore) SaveReminderState(
	context.Context,
	*points.ReminderState,
) error {
	m.saveReminderCalled = true
	return nil
}

func TestManagerConsume(t *testing.T) {
	store := &mockStore{}

	manager := &Manager{
		store: store,
		pointSystems: map[string]*points.PointSystem{
			"TP": {
				Definition: points.Definition{
					ID:            "TP",
					Name:          "Training Points",
					Max:           100,
					RegenMinutes: 10,
				},
				Current: 50,
				Elapsed: 5 * time.Minute,
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

	if err := manager.Consume(context.Background(), "TP", 10); err != nil {
		t.Fatalf("Consume() error = %v", err)
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

func TestManagerConsume_AddPoints(t *testing.T) {
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
				Current: 40,
				Elapsed: 5 * time.Minute,
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

	if err := manager.Consume(context.Background(), "TP", -25); err != nil {
		t.Fatalf("Consume() error = %v", err)
	}

	system := manager.pointSystems["TP"]

	if system.Current != 65 {
		t.Fatalf("Current = %d, want 65", system.Current)
	}

	if system.Elapsed != 5*time.Minute {
		t.Fatalf("Elapsed = %v, want 5m", system.Elapsed)
	}

	reminder := manager.reminders["TP"]

	if reminder.AlertSent {
		t.Fatal("AlertSent should be reset")
	}

	if reminder.FullSent {
		t.Fatal("FullSent should be reset")
	}
}

func TestManagerConsume_UnknownSystem(t *testing.T) {
	manager := &Manager{
		store: &mockStore{},
	}

	if err := manager.Consume(context.Background(), "UNKNOWN", 10); err == nil {
		t.Fatal("expected error")
	}
}

func TestManagerSet(t *testing.T) {
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
				Current: 40,
				Elapsed: 5 * time.Minute,
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

	if err := manager.Set(context.Background(), "TP", 75); err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	system := manager.pointSystems["TP"]

	if system.Current != 75 {
		t.Fatalf("Current = %d, want 75", system.Current)
	}

	if system.Elapsed != 0 {
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

func TestManagerSet_UnknownSystem(t *testing.T) {
	manager := &Manager{
		store: &mockStore{},
	}

	if err := manager.Set(context.Background(), "UNKNOWN", 10); err == nil {
		t.Fatal("expected error")
	}
}
