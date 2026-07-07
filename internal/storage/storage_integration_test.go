package storage

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"umamusume-notifier/internal/points"
)

func TestSQLiteStore_RoundTrip(t *testing.T) {
	ctx := context.Background()

	dbPath := filepath.Join(t.TempDir(), "test.db")

	store, err := NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore() error = %v", err)
	}
	defer store.Close()

	if err := store.Initialize(ctx); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	definitions := []points.Definition{
		{
			ID:           "TP",
			Name:         "Training Points",
			Max:          100,
			RegenMinutes: 10,
		},
		{
			ID:           "RP",
			Name:         "Race Points",
			Max:          5,
			RegenMinutes: 120,
		},
	}

	if err := store.SyncPointSystems(ctx, definitions); err != nil {
		t.Fatalf("SyncPointSystems() error = %v", err)
	}

	systems, err := store.LoadPointSystems(ctx)
	if err != nil {
		t.Fatalf("LoadPointSystems() error = %v", err)
	}

	if len(systems) != 2 {
		t.Fatalf("expected 2 systems, got %d", len(systems))
	}

	systems[0].Current = 42
	systems[0].Elapsed = 7 * time.Minute
	systems[0].LastTick = time.Now().UTC()

	if err := store.SavePointSystems(ctx, systems); err != nil {
		t.Fatalf("SavePointSystems() error = %v", err)
	}

	reloaded, err := store.LoadPointSystems(ctx)
	if err != nil {
		t.Fatalf("Reload error = %v", err)
	}

	if reloaded[0].Current != 42 {
		t.Errorf("Current = %d, want 42", reloaded[0].Current)
	}

	if reloaded[0].Elapsed != 7*time.Minute {
		t.Errorf("Elapsed = %v, want %v",
			reloaded[0].Elapsed,
			7*time.Minute)
	}

	state := &points.ReminderState{
		SystemID:      "TP",
		AlertSent:     true,
		FullSent:      false,
		LastMessageID: 12345,
	}

	if err := store.SaveReminderState(ctx, state); err != nil {
		t.Fatalf("SaveReminderState() error = %v", err)
	}

	states, err := store.LoadReminderStates(ctx)
	if err != nil {
		t.Fatalf("LoadReminderStates() error = %v", err)
	}

	if len(states) != 1 {
		t.Fatalf("expected 1 reminder state, got %d", len(states))
	}

	if states[0].SystemID != "TP" {
		t.Errorf("SystemID = %q, want TP", states[0].SystemID)
	}

	if !states[0].AlertSent {
		t.Error("AlertSent = false, want true")
	}

	if states[0].LastMessageID != 12345 {
		t.Errorf("LastMessageID = %d, want 12345", states[0].LastMessageID)
	}

	if len(systems) != len(definitions) {
		t.Fatalf("expected %d systems, got %d", len(definitions), len(systems))
	}

	if systems[0].ID != "RP" {
		t.Errorf("systems[0].ID = %q, want RP", systems[0].ID)
	}

	if systems[1].ID != "TP" {
		t.Errorf("systems[1].ID = %q, want TP", systems[1].ID)
	}

	if systems[1].Max != 100 {
		t.Errorf("Max = %d, want 100", systems[1].Max)
	}

	if systems[1].Current != 0 {
		t.Errorf("Current = %d, want 0", systems[1].Current)
	}

	loaded := make(map[string]*points.PointSystem)

	for _, ps := range systems {
		loaded[ps.ID] = ps
	}

	tp := loaded["TP"]

	if tp == nil {
		t.Fatal("TP not loaded")
	}

	if tp.Max != 100 {
		t.Errorf("Max = %d, want 100", tp.Max)
	}

	if tp.Current != 0 {
		t.Errorf("Current = %d, want 0", tp.Current)
	}
}
