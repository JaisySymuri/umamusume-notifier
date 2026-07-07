package notification

type Type int

const (
	NearFull Type = iota
	Full
)

type Event struct {
	SystemID string
	SystemName string
	Type     Type
}
