package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"umamusume-notifier/internal/metrics"
	"umamusume-notifier/internal/notification"
	"umamusume-notifier/internal/points"
)

type tickMockStore struct{}

func (m *tickMockStore) Initialize(context.Context) error                            { return nil }
func (m *tickMockStore) SyncPointSystems(context.Context, []points.Definition) error { return nil }
func (m *tickMockStore) LoadPointSystems(context.Context) ([]*points.PointSystem, error) {
	return nil, nil
}
func (m *tickMockStore) LoadReminderStates(context.Context) ([]*points.ReminderState, error) {
	return nil, nil
}
func (m *tickMockStore) SavePointSystems(context.Context, []*points.PointSystem) error  { return nil }
func (m *tickMockStore) SaveReminderState(context.Context, *points.ReminderState) error { return nil }

func TestManagerTickRecordsScheduledMetric(t *testing.T) {
	manager := &Manager{
		store: &tickMockStore{},
		pointSystems: map[string]*points.PointSystem{
			"TP": {
				Definition: points.Definition{
					ID:           "TP",
					Name:         "Training Points",
					Max:          100,
					RegenMinutes: 10,
				},
				Current:  100,
				Elapsed:  0,
				LastTick: time.Now().Add(-10 * time.Minute),
			},
		},
		reminders: map[string]*points.ReminderState{
			"TP": {
				SystemID:  "TP",
				FullSince: time.Now().Add(-2 * time.Hour),
			},
		},
		alertThreshold: 30 * time.Minute,
	}

	events, err := manager.Tick(context.Background(), time.Now())
	if err != nil {
		t.Fatalf("Tick() error = %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("got %d events, want 1", len(events))
	}

	if events[0].Type != notification.Full {
		t.Fatalf("event type = %v, want Full", events[0].Type)
	}

	if events[0].ScheduledAt.IsZero() {
		t.Fatal("ScheduledAt should be set")
	}

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()
	metrics.Handler().ServeHTTP(rr, req)

	if !strings.Contains(rr.Body.String(), "bot_reminders_scheduled_total 1") {
		t.Fatalf("expected scheduled metric, got:\n%s", rr.Body.String())
	}

	if !strings.Contains(rr.Body.String(), "bot_full_over_hour_total 1") {
		t.Fatalf("expected over-hour metric, got:\n%s", rr.Body.String())
	}
}
