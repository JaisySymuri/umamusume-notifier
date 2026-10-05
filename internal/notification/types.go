package notification

import "time"

type Type int

const (
	NearFull Type = iota
	Full
	Milestone
)

type Event struct {
	SystemID    string
	SystemName  string
	Type        Type
	Points      int
	ScheduledAt time.Time
}
