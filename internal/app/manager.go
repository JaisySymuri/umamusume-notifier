package app

import (
	"sync"
	"time"

	"umamusume-notifier/internal/points"
	"umamusume-notifier/internal/storage"
)

type Manager struct {
    mu sync.RWMutex

    // Dependencies
    store storage.Store

    // State
    pointSystems map[string]*points.PointSystem
    reminders    map[string]*points.ReminderState

    // Configuration
    alertThreshold time.Duration
}

func New(
	store storage.Store,
	alertThreshold time.Duration,
) *Manager {

	return &Manager{
		store:          store,
		pointSystems: make(map[string]*points.PointSystem),
		reminders:      make(map[string]*points.ReminderState),
		alertThreshold: alertThreshold,
	}
}

func (m *Manager) system(id string) (
	*points.PointSystem,
	*points.ReminderState,
	bool,
) {
	system, ok := m.pointSystems[id]
	if !ok {
		return nil, nil, false
	}

	reminder, ok := m.reminders[id]
	if !ok {
		return nil, nil, false
	}

	return system, reminder, true
}