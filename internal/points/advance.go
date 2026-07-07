package points

import "time"

// regenDuration returns the duration required to regenerate one point.
func (p *PointSystem) regenDuration() time.Duration {
	return time.Duration(p.RegenMinutes) * time.Minute
}

// Advance moves the system forward by delta and regenerates points whenever
// enough elapsed time has accumulated.
//
// It is safe to call with any duration, including durations spanning multiple
// regeneration intervals.
func (p *PointSystem) Advance(delta time.Duration) {
	if delta <= 0 {
		return
	}

	if p.Current >= p.Max {
		p.Current = p.Max
		p.Elapsed = 0
		return
	}

	p.Elapsed += delta

	regen := p.regenDuration()

	for p.Elapsed >= regen && p.Current < p.Max {
		p.Elapsed -= regen
		p.Current++
	}

	if p.Current >= p.Max {
		p.Current = p.Max
		p.Elapsed = 0
	}
}