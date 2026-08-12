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

func TestFormatStatus_Grouped(t *testing.T) {
	statuses := []app.Status{
		{
			ID:            "CP",
			Name:          "Club Points",
			Current:       0,
			Max:           1,
			TimeUntilFull: 17 * time.Minute,
		},
		{
			ID:            "RP",
			Name:          "Race Points",
			Current:       4,
			Max:           5,
			TimeUntilFull: 28 * time.Minute,
		},
		{
			ID:      "TP",
			Name:    "Training Points",
			Current: 100,
			Max:     100,
			Full:    true,
		},
		{
			ID:            "1P",
			Name:          "1st Island Points",
			Current:       0,
			Max:           1,
			TimeUntilFull: 8 * time.Hour,
		},
		{
			ID:            "2P",
			Name:          "2nd Island Points",
			Current:       0,
			Max:           1,
			TimeUntilFull: 10 * time.Hour,
		},
		{
			ID:            "LP",
			Name:          "Live Points",
			Current:       0,
			Max:           20,
			TimeUntilFull: 15 * time.Hour,
		},
		{
			ID:            "MP",
			Name:          "Minigame Points",
			Current:       0,
			Max:           200,
			TimeUntilFull: 23*time.Hour + 20*time.Minute,
		},
	}

	got := FormatStatus(statuses)

	assertContains := func(substr string) {
		t.Helper()

		if !strings.Contains(got, substr) {
			t.Fatalf("output does not contain %q\n\ngot:\n%s", substr, got)
		}
	}

	assertContains("Umamusume")
	assertContains("Club Points (CP)")
	assertContains("Race Points (RP)")
	assertContains("Training Points (TP)")
	assertContains("Holodori")
	assertContains("1st Island Points (1P)")
	assertContains("2nd Island Points (2P)")
	assertContains("Live Points (LP)")
	assertContains("Minigame Points (MP)")
	assertContains("  FULL")
	assertContains("  Full in: 17m (")
	assertContains("  Full in: 23h 20m (")
}

func TestFormatStatus_UnknownGroup(t *testing.T) {
	got := FormatStatus([]app.Status{
		{
			ID:            "XX",
			Name:          "Mystery Points",
			Current:       2,
			Max:           3,
			TimeUntilFull: 30 * time.Minute,
		},
	})

	if !strings.Contains(got, "Mystery Points (XX)") {
		t.Fatalf("unexpected output:\n%s", got)
	}

	if strings.Contains(got, "Umamusume") || strings.Contains(got, "Holodori") {
		t.Fatalf("unexpected known group in output:\n%s", got)
	}
}
