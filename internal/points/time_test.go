package points

import (
	"testing"
	"time"
)

func TestIsFull(t *testing.T) {
	tests := []struct {
		name     string
		current  int
		max      int
		wantFull bool
	}{
		{
			name:     "empty",
			current:  0,
			max:      5,
			wantFull: false,
		},
		{
			name:     "partial",
			current:  3,
			max:      5,
			wantFull: false,
		},
		{
			name:     "full",
			current:  5,
			max:      5,
			wantFull: true,
		},
		{
			name:     "greater than max",
			current:  6,
			max:      5,
			wantFull: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newTestPointSystem(
				tt.current,
				tt.max,
				10,
				0,
			)

			if got := p.IsFull(); got != tt.wantFull {
				t.Errorf("IsFull() = %v, want %v", got, tt.wantFull)
			}
		})
	}
}

func TestTimeUntilFull(t *testing.T) {
	tests := []struct {
		name     string
		current  int
		max      int
		regen    int
		elapsed  time.Duration
		wantTime time.Duration
	}{
		{
			name:     "already full",
			current:  5,
			max:      5,
			regen:    10,
			elapsed:  0,
			wantTime: 0,
		},
		{
			name:     "one point remaining",
			current:  4,
			max:      5,
			regen:    10,
			elapsed:  0,
			wantTime: 10 * time.Minute,
		},
		{
			name:     "partial elapsed",
			current:  4,
			max:      5,
			regen:    10,
			elapsed:  3 * time.Minute,
			wantTime: 7 * time.Minute,
		},
		{
			name:     "multiple points remaining",
			current:  2,
			max:      5,
			regen:    10,
			elapsed:  4 * time.Minute,
			wantTime: 26 * time.Minute,
		},
		{
			name:     "empty system",
			current:  0,
			max:      5,
			regen:    10,
			elapsed:  0,
			wantTime: 50 * time.Minute,
		},
		{
			name:     "elapsed exceeds remaining",
			current:  4,
			max:      5,
			regen:    10,
			elapsed:  15 * time.Minute,
			wantTime: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newTestPointSystem(
				tt.current,
				tt.max,
				tt.regen,
				tt.elapsed,
			)

			if got := p.TimeUntilFull(); got != tt.wantTime {
				t.Errorf("TimeUntilFull() = %v, want %v", got, tt.wantTime)
			}
		})
	}
}