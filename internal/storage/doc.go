// Package storage provides persistence for the application's runtime state.
//
// The storage layer owns mutable state such as current points, elapsed
// regeneration time, reminder flags, and last update timestamps.
//
// Point system definitions are owned by the config package and synchronized
// into storage during application startup.
package storage