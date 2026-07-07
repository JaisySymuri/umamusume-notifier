package points

import (
	"testing"
	"time"
)

func TestSetElapsed(t *testing.T) {
	tests := []struct {
		name        string
		current     int
		max         int
		regen       int
		initial     time.Duration
		elapsed     time.Duration
		wantCurrent int
		wantElapsed time.Duration
	}{
		{
			name:        "less than one regen",
			current:     2,
			max:         5,
			regen:       10,
			initial:     5 * time.Minute,
			elapsed:     7 * time.Minute,
			wantCurrent: 2,
			wantElapsed: 7 * time.Minute,
		},
		{
			name:        "exactly one regen",
			current:     2,
			max:         5,
			regen:       10,
			initial:     3 * time.Minute,
			elapsed:     10 * time.Minute,
			wantCurrent: 3,
			wantElapsed: 0,
		},
		{
			name:        "overflow one point",
			current:     2,
			max:         5,
			regen:       10,
			initial:     1 * time.Minute,
			elapsed:     17 * time.Minute,
			wantCurrent: 3,
			wantElapsed: 7 * time.Minute,
		},
		{
			name:        "overflow several points",
			current:     1,
			max:         5,
			regen:       10,
			initial:     2 * time.Minute,
			elapsed:     35 * time.Minute,
			wantCurrent: 4,
			wantElapsed: 5 * time.Minute,
		},
		{
			name:        "reach max",
			current:     3,
			max:         5,
			regen:       10,
			initial:     4 * time.Minute,
			elapsed:     40 * time.Minute,
			wantCurrent: 5,
			wantElapsed: 0,
		},
		{
			name:        "already full",
			current:     5,
			max:         5,
			regen:       10,
			initial:     8 * time.Minute,
			elapsed:     5 * time.Minute,
			wantCurrent: 5,
			wantElapsed: 0,
		},
		{
			name:        "zero elapsed",
			current:     2,
			max:         5,
			regen:       10,
			initial:     6 * time.Minute,
			elapsed:     0,
			wantCurrent: 2,
			wantElapsed: 0,
		},
		{
			name:        "negative elapsed ignored",
			current:     2,
			max:         5,
			regen:       10,
			initial:     6 * time.Minute,
			elapsed:     -5 * time.Minute,
			wantCurrent: 2,
			wantElapsed: 6 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PointSystem{
				Definition: Definition{
					ID:           testID,
					Name:         testName,
					Max:          tt.max,
					RegenMinutes: tt.regen,
				},
				Current: tt.current,
				Elapsed: tt.initial,
			}

			p.SetElapsed(tt.elapsed)

			if p.Current != tt.wantCurrent {
				t.Errorf("Current = %d, want %d", p.Current, tt.wantCurrent)
			}

			if p.Elapsed != tt.wantElapsed {
				t.Errorf("Elapsed = %v, want %v", p.Elapsed, tt.wantElapsed)
			}
		})
	}
}