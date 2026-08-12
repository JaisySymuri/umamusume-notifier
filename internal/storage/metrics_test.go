package storage

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"umamusume-notifier/internal/metrics"
	"umamusume-notifier/internal/points"
)

func TestSQLiteStore_RecordsStorageMetrics(t *testing.T) {
	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "metrics.db")

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
	}

	if err := store.SyncPointSystems(ctx, definitions); err != nil {
		t.Fatalf("SyncPointSystems() error = %v", err)
	}

	if _, err := store.LoadPointSystems(ctx); err != nil {
		t.Fatalf("LoadPointSystems() error = %v", err)
	}

	if _, err := store.LoadReminderStates(ctx); err != nil {
		t.Fatalf("LoadReminderStates() error = %v", err)
	}

	if err := store.SavePointSystems(ctx, []*points.PointSystem{{
		Definition: definitions[0],
	}}); err != nil {
		t.Fatalf("SavePointSystems() error = %v", err)
	}

	if err := store.SaveReminderState(ctx, &points.ReminderState{SystemID: "TP"}); err != nil {
		t.Fatalf("SaveReminderState() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()
	metrics.Handler().ServeHTTP(rr, req)

	body := rr.Body.String()
	for _, want := range []string{
		`bot_storage_op_duration_seconds_bucket`,
		`op="initialize"`,
		`op="sync_point_systems"`,
		`op="load_point_systems"`,
		`op="load_reminder_states"`,
		`op="save_point_systems"`,
		`op="save_reminder_state"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("metrics output does not contain %q\nbody:\n%s", want, body)
		}
	}
}
