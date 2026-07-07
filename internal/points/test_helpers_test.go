package points

import "time"

const (
	testID   = "TP"
	testName = "Training Points"
)

// newTestPointSystem creates a PointSystem with sensible defaults for testing.
func newTestPointSystem(current, max, regenMinutes int, elapsed time.Duration) *PointSystem {
	return &PointSystem{
		Definition: Definition{
			ID:           testID,
			Name:         testName,
			Max:          max,
			RegenMinutes: regenMinutes,
		},
		Current: current,
		Elapsed: elapsed,
		LastTick: time.Now().UTC(),
	}
}