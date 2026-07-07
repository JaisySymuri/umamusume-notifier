package points

import "errors"

var (
	ErrEmptyID          = errors.New("points: id cannot be empty")
	ErrEmptyName        = errors.New("points: name cannot be empty")
	ErrInvalidMax       = errors.New("points: max must be greater than zero")
	ErrInvalidRegenTime = errors.New("points: regen minutes must be greater than zero")
)