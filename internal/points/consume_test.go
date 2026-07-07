package points

import (
	"testing"
	"time"
)

func TestConsume(t *testing.T) {
	tests := []struct {
		name        string
		current     int
		max         int
		elapsed     time.Duration
		amount      int
		wantCurrent int
		wantElapsed time.Duration
	}{
		{
			name:        "consume one",
			current:     5,
			max:         10,
			elapsed:     7 * time.Minute,
			amount:      1,
			wantCurrent: 4,
			wantElapsed: 0,
		},
		{
			name:        "consume many",
			current:     8,
			max:         10,
			elapsed:     5 * time.Minute,
			amount:      3,
			wantCurrent: 5,
			wantElapsed: 0,
		},
		{
			name:        "consume all",
			current:     5,
			max:         10,
			elapsed:     8 * time.Minute,
			amount:      5,
			wantCurrent: 0,
			wantElapsed: 0,
		},
		{
			name:        "consume more than current",
			current:     2,
			max:         10,
			elapsed:     4 * time.Minute,
			amount:      10,
			wantCurrent: 0,
			wantElapsed: 0,
		},
		{
			name:        "consume zero ignored",
			current:     5,
			max:         10,
			elapsed:     3 * time.Minute,
			amount:      0,
			wantCurrent: 5,
			wantElapsed: 3 * time.Minute,
		},
		{
			name:        "negative consume ignored",
			current:     5,
			max:         10,
			elapsed:     3 * time.Minute,
			amount:      -5,
			wantCurrent: 5,
			wantElapsed: 3 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PointSystem{
				Definition: Definition{
					ID:           "TP",
					Name:         "Training Points",
					Max:          tt.max,
					RegenMinutes: 10,
				},
				Current:      tt.current,
				Elapsed:      tt.elapsed,
			}

			p.Consume(tt.amount)

			if p.Current != tt.wantCurrent {
				t.Errorf("Current = %d, want %d", p.Current, tt.wantCurrent)
			}

			if p.Elapsed != tt.wantElapsed {
				t.Errorf("Elapsed = %v, want %v", p.Elapsed, tt.wantElapsed)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name        string
		current     int
		max         int
		elapsed     time.Duration
		amount      int
		wantCurrent int
		wantElapsed time.Duration
	}{
		{
			name:        "add one",
			current:     5,
			max:         10,
			elapsed:     4 * time.Minute,
			amount:      1,
			wantCurrent: 6,
			wantElapsed: 4 * time.Minute,
		},
		{
			name:        "add many",
			current:     2,
			max:         10,
			elapsed:     5 * time.Minute,
			amount:      4,
			wantCurrent: 6,
			wantElapsed: 5 * time.Minute,
		},
		{
			name:        "add to max resets elapsed",
			current:     8,
			max:         10,
			elapsed:     9 * time.Minute,
			amount:      2,
			wantCurrent: 10,
			wantElapsed: 0,
		},
		{
			name:        "add beyond max",
			current:     9,
			max:         10,
			elapsed:     8 * time.Minute,
			amount:      10,
			wantCurrent: 10,
			wantElapsed: 0,
		},
		{
			name:        "add zero ignored",
			current:     5,
			max:         10,
			elapsed:     6 * time.Minute,
			amount:      0,
			wantCurrent: 5,
			wantElapsed: 6 * time.Minute,
		},
		{
			name:        "negative add ignored",
			current:     5,
			max:         10,
			elapsed:     6 * time.Minute,
			amount:      -2,
			wantCurrent: 5,
			wantElapsed: 6 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PointSystem{
				Definition: Definition{
					ID:           "TP",
					Name:         "Training Points",
					Max:          tt.max,
					RegenMinutes: 10,
				},
				Current:      tt.current,
				Elapsed:      tt.elapsed,
			}

			p.Add(tt.amount)

			if p.Current != tt.wantCurrent {
				t.Errorf("Current = %d, want %d", p.Current, tt.wantCurrent)
			}

			if p.Elapsed != tt.wantElapsed {
				t.Errorf("Elapsed = %v, want %v", p.Elapsed, tt.wantElapsed)
			}
		})
	}
}