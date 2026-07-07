package points

import "time"

// IsFull reports whether the point system has reached its maximum capacity.
func (p PointSystem) IsFull() bool {
	return p.Current >= p.Max
}

// TimeUntilFull returns the remaining time required for the point system to
// reach its maximum capacity.
//
// If the system is already full, zero is returned.
func (p PointSystem) TimeUntilFull() time.Duration {
	if p.IsFull() {
		return 0
	}

	remainingPoints := p.Max - p.Current
	remaining := p.regenDuration()*time.Duration(remainingPoints) - p.Elapsed

	if remaining < 0 {
		return 0
	}

	return remaining
}