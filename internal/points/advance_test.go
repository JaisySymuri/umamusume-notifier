package points

import (
	"testing"
	"time"
)

func TestAdvance(t *testing.T) {
	tests := []struct {
		name         string
		current      int
		max          int
		regenMinutes int
		elapsed      time.Duration
		delta        time.Duration
		wantCurrent  int
		wantElapsed  time.Duration
	}{
		{
			name:         "no regen yet",
			current:      0,
			max:          5,
			regenMinutes: 10,
			elapsed:      0,
			delta:        5 * time.Minute,
			wantCurrent:  0,
			wantElapsed:  5 * time.Minute,
		},
		{
			name:         "exactly one regen",
			current:      0,
			max:          5,
			regenMinutes: 10,
			elapsed:      0,
			delta:        10 * time.Minute,
			wantCurrent:  1,
			wantElapsed:  0,
		},
		{
			name:         "carry remaining elapsed",
			current:      2,
			max:          5,
			regenMinutes: 10,
			elapsed:      7 * time.Minute,
			delta:        5 * time.Minute,
			wantCurrent:  3,
			wantElapsed:  2 * time.Minute,
		},
		{
			name:         "multiple regen intervals",
			current:      1,
			max:          5,
			regenMinutes: 10,
			elapsed:      0,
			delta:        35 * time.Minute,
			wantCurrent:  4,
			wantElapsed:  5 * time.Minute,
		},
		{
			name:         "reach max and reset elapsed",
			current:      3,
			max:          5,
			regenMinutes: 10,
			elapsed:      7 * time.Minute,
			delta:        40 * time.Minute,
			wantCurrent:  5,
			wantElapsed:  0,
		},
		{
			name:         "already full",
			current:      5,
			max:          5,
			regenMinutes: 10,
			elapsed:      5 * time.Minute,
			delta:        30 * time.Minute,
			wantCurrent:  5,
			wantElapsed:  0,
		},
		{
			name:         "zero delta",
			current:      2,
			max:          5,
			regenMinutes: 10,
			elapsed:      3 * time.Minute,
			delta:        0,
			wantCurrent:  2,
			wantElapsed:  3 * time.Minute,
		},
		{
			name:         "negative delta ignored",
			current:      2,
			max:          5,
			regenMinutes: 10,
			elapsed:      3 * time.Minute,
			delta:        -5 * time.Minute,
			wantCurrent:  2,
			wantElapsed:  3 * time.Minute,
		},
		{
			name:         "huge delta",
			current:      0,
			max:          5,
			regenMinutes: 10,
			elapsed:      0,
			delta:        24 * time.Hour,
			wantCurrent:  5,
			wantElapsed:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PointSystem{
				Definition: Definition{
					ID:           "TP",
					Name:         "Training Points",
					Max:          tt.max,
					RegenMinutes: tt.regenMinutes,
				},
				Current:      tt.current,
				Elapsed:      tt.elapsed,
			}

			p.Advance(tt.delta)

			if p.Current != tt.wantCurrent {
				t.Errorf("Current = %d, want %d", p.Current, tt.wantCurrent)
			}

			if p.Elapsed != tt.wantElapsed {
				t.Errorf("Elapsed = %v, want %v", p.Elapsed, tt.wantElapsed)
			}
		})
	}
}

