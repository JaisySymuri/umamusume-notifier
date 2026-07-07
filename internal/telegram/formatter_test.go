package telegram

import (
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

	const want = `📊 Point Status

Combat Points (CP)
  1/1
  FULL

Training Points (TP)
  80/100
  Full in: 2h 30m
`

	if got != want {
		t.Fatalf("got:\n%s\n\nwant:\n%s", got, want)
	}
}
