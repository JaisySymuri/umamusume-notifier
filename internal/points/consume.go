package points

// Consume spends points and resets the regeneration progress.
//
// Consuming more points than currently available will clamp the current
// points to zero.
func (p *PointSystem) Consume(amount int) {
	if amount <= 0 {
		return
	}

	p.Current -= amount

	if p.Current < 0 {
		p.Current = 0
	}

	// Consuming always restarts the regeneration timer.
	p.Elapsed = 0
}

// Add grants points up to the configured maximum.
//
// Adding points does not affect the current regeneration progress.
func (p *PointSystem) Add(amount int) {
	if amount <= 0 {
		return
	}

	p.Current += amount

	if p.Current > p.Max {
		p.Current = p.Max
	}

	if p.Current == p.Max {
		p.Elapsed = 0
	}
}

// Set assigns the current points directly.
//
// Values below zero are clamped to zero and values above the maximum are
// clamped to the configured maximum. Setting the current points resets the
// regeneration progress because the timer no longer matches the visible state.
func (p *PointSystem) Set(amount int) {
	p.Current = amount

	if p.Current < 0 {
		p.Current = 0
	}

	if p.Current > p.Max {
		p.Current = p.Max
	}

	p.Elapsed = 0
}
