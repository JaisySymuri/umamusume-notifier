package points

import "time"

// SetElapsed replaces the current regeneration progress.
//
// If the supplied elapsed duration exceeds one or more regeneration intervals,
// the excess is automatically converted into regenerated points by reusing
// Advance(). Negative durations are ignored.
func (p *PointSystem) SetElapsed(elapsed time.Duration) {
	if elapsed < 0 {
		return
	}

	p.Elapsed = 0
	p.Advance(elapsed)
}