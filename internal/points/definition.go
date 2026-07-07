package points

// Definition describes the immutable configuration of a point system.
type Definition struct {
	ID string

	Name string

	Max int

	// RegenMinutes is the number of minutes required to regenerate one point.
	RegenMinutes int
}