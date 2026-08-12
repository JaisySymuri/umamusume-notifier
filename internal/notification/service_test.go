package notification

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"umamusume-notifier/internal/metrics"
)

type mockSender struct {
	messageID int
	err       error
}

func (m *mockSender) Notify(int64, Event) (int, error) {
	return m.messageID, m.err
}

type mockRecorder struct {
	called    bool
	systemID  string
	messageID int
	err       error
}

func (m *mockRecorder) RecordReminderMessage(
	_ context.Context,
	systemID string,
	messageID int,
) error {
	m.called = true
	m.systemID = systemID
	m.messageID = messageID
	return m.err
}

func metricsBody(t *testing.T) string {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()

	metrics.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("metrics status = %d, want %d", rr.Code, http.StatusOK)
	}

	return rr.Body.String()
}

func TestServiceNotifySuccessRecordsMetrics(t *testing.T) {
	sender := &mockSender{messageID: 42}
	recorder := &mockRecorder{}
	service := NewService(sender, recorder, 123)

	if err := service.Notify(context.Background(), Event{
		SystemID:    "TP",
		SystemName:  "Training Points",
		Type:        NearFull,
		ScheduledAt: time.Now().Add(-90 * time.Second),
	}); err != nil {
		t.Fatalf("Notify() error = %v", err)
	}

	if !recorder.called {
		t.Fatal("expected reminder to be recorded")
	}

	body := metricsBody(t)
	if !strings.Contains(body, `bot_reminders_sent_total{outcome="success"} 1`) {
		t.Fatalf("expected success metric, got:\n%s", body)
	}

	if !strings.Contains(body, "bot_reminder_delivery_delay_seconds_sum") {
		t.Fatalf("expected delivery delay histogram, got:\n%s", body)
	}
}

func TestServiceNotifyErrorRecordsMetrics(t *testing.T) {
	sender := &mockSender{err: errors.New("boom")}
	recorder := &mockRecorder{}
	service := NewService(sender, recorder, 123)

	if err := service.Notify(context.Background(), Event{
		SystemID:   "TP",
		SystemName: "Training Points",
		Type:       NearFull,
	}); err == nil {
		t.Fatal("expected error")
	}

	if recorder.called {
		t.Fatal("reminder should not be recorded on send failure")
	}

	body := metricsBody(t)
	if !strings.Contains(body, `bot_reminders_sent_total{outcome="error"} 1`) {
		t.Fatalf("expected error metric, got:\n%s", body)
	}
}
