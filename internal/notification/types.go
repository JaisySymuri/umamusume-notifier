package notification

import "time"

type Type int

const (
	NearFull Type = iota
	Full
)

type Event struct {
	SystemID    string
	SystemName  string
	Type        Type
	ScheduledAt time.Time
}
