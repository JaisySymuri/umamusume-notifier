package telegram

import (
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"umamusume-notifier/internal/metrics"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type fakeTelegramAPI struct {
	messageID int
	err       error
}

func (f fakeTelegramAPI) Send(tgbotapi.Chattable) (tgbotapi.Message, error) {
	if f.err != nil {
		return tgbotapi.Message{}, f.err
	}

	return tgbotapi.Message{MessageID: f.messageID}, nil
}

func (f fakeTelegramAPI) GetUpdatesChan(tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel {
	return nil
}

type timeoutErr struct{}

func (timeoutErr) Error() string   { return "timeout" }
func (timeoutErr) Timeout() bool   { return true }
func (timeoutErr) Temporary() bool { return false }

func TestClassifyTelegramError(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want string
	}{
		{name: "timeout", err: timeoutErr{}, want: "timeout"},
		{name: "network", err: &net.DNSError{IsTimeout: false, IsTemporary: true}, want: "network"},
		{name: "rate limit", err: errors.New("Too Many Requests: retry after 10"), want: "rate_limit"},
		{name: "telegram api", err: errors.New("telegram api response: bad request"), want: "telegram_api"},
		{name: "unknown", err: errors.New("something else"), want: "unknown"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := classifyTelegramError(tc.err); got != tc.want {
				t.Fatalf("classifyTelegramError() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestTelegramSenderSendSuccess(t *testing.T) {
	sender := &telegramSender{api: fakeTelegramAPI{messageID: 77}}

	got, err := sender.Send(123, "hello")
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if got != 77 {
		t.Fatalf("Send() = %d, want 77", got)
	}
}

func TestTelegramSenderSendError(t *testing.T) {
	sender := &telegramSender{api: fakeTelegramAPI{err: errors.New("boom")}}

	_, err := sender.Send(123, "hello")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestMetricsExposeTelegramAPI(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()

	metrics.Handler().ServeHTTP(rr, req)
	body := rr.Body.String()

	for _, want := range []string{
		"bot_telegram_api_requests_total",
		"bot_telegram_api_duration_seconds",
		"bot_telegram_api_errors_total",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("metrics output does not contain %q\nbody:\n%s", want, body)
		}
	}
}
