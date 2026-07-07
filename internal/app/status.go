package app

import "time"
import "sort"


// Status represents the information needed to display a point system.
type Status struct {
	ID            string
	Name          string
	Current       int
	Max           int
	TimeUntilFull time.Duration
	Full          bool
}

func (m *Manager) Status() []Status {
	m.mu.RLock()
	defer m.mu.RUnlock()

	statuses := make([]Status, 0, len(m.pointSystems))

	for _, system := range m.pointSystems {
		statuses = append(statuses, Status{
			ID:            system.ID,
			Name:          system.Name,
			Current:       system.Current,
			Max:           system.Max,
			TimeUntilFull: system.TimeUntilFull(),
			Full:          system.Current >= system.Max,
		})
	}

	sort.Slice(statuses, func(i, j int) bool {
		return statuses[i].ID < statuses[j].ID
	})

	return statuses
}