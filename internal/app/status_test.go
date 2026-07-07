package app

import (
	"testing"
	"time"

	"umamusume-notifier/internal/points"
)

func TestManagerStatus_Empty(t *testing.T) {
	manager := &Manager{}

	got := manager.Status()

	if len(got) != 0 {
		t.Fatalf("got %d statuses, want 0", len(got))
	}
}

func TestManagerStatus(t *testing.T) {
	manager := &Manager{
		pointSystems: map[string]*points.PointSystem{
			"TP": {
				Definition: points.Definition{
					ID:            "TP",
					Name:          "Training Points",
					Max:           100,
					RegenMinutes: 10,
				},
				Current: 80,
				Elapsed: 0,
			},
			"CP": {
				Definition: points.Definition{
					ID:            "CP",
					Name:          "Combat Points",
					Max:           1,
					RegenMinutes: 480,
				},
				Current: 1,
			},
		},
	}

	got := manager.Status()

	if len(got) != 2 {
		t.Fatalf("got %d systems, want 2", len(got))
	}

	if got[0].ID != "CP" {
		t.Fatalf("first system = %q, want CP", got[0].ID)
	}

	if got[1].ID != "TP" {
		t.Fatalf("second system = %q, want TP", got[1].ID)
	}

	if !got[0].Full {
		t.Fatal("CP should be full")
	}

	if got[1].Full {
		t.Fatal("TP should not be full")
	}

	want := 200 * time.Minute

	if got[1].TimeUntilFull != want {
		t.Fatalf(
			"time until full = %v, want %v",
			got[1].TimeUntilFull,
			want,
		)
	}
}