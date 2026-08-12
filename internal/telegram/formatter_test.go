package telegram

import (
	"strings"
	"testing"
	"time"

	"umamusume-notifier/internal/app"
)

func TestFormatStatus_Empty(t *testing.T) {
	got := FormatStatus(nil)

	want := "No point systems configured."

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestFormatStatus(t *testing.T) {
	statuses := []app.Status{
		{
			ID:      "CP",
			Name:    "Combat Points",
			Current: 1,
			Max:     1,
			Full:    true,
		},
		{
			ID:            "TP",
			Name:          "Training Points",
			Current:       80,
			Max:           100,
			TimeUntilFull: 2*time.Hour + 30*time.Minute,
		},
	}

	got := FormatStatus(statuses)

	if got == "" {
		t.Fatal("FormatStatus() returned empty output")
	}

	assertContains := func(substr string) {
		t.Helper()

		if !strings.Contains(got, substr) {
			t.Fatalf("output does not contain %q\n\ngot:\n%s", substr, got)
		}
	}

	assertContains("Point Status")
	assertContains("Combat Points (CP)")
	assertContains("  1/1")
	assertContains("  FULL")
	assertContains("Training Points (TP)")
	assertContains("  80/100")
	assertContains("  Full in: 2h 30m (")
	assertContains(" WIB)")
}
