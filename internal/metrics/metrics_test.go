package metrics

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListenAddr(t *testing.T) {
	if got := ListenAddr(); got != "127.0.0.1:9091" {
		t.Fatalf("ListenAddr() = %q, want %q", got, "127.0.0.1:9091")
	}
}

func TestHandlerExposesMetrics(t *testing.T) {
	ObserveCommand("/status", "success")

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()

	Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", rr.Code, http.StatusOK)
	}

	body := rr.Body.String()
	for _, want := range []string{
		"bot_commands_total",
		"go_goroutines",
		"go_memstats_alloc_bytes",
		"process_cpu_seconds_total",
		"process_resident_memory_bytes",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("metrics output does not contain %q\nbody:\n%s", want, body)
		}
	}
}
