package app

import (
	"context"
	"fmt"
)

func (m *Manager) ConsumeReply(
	ctx context.Context,
	messageID int,
	amount int,
) error {
	systemID, ok := m.systemIDByMessageID(messageID)
	if !ok {
		return fmt.Errorf("unknown reminder message %d", messageID)
	}

	return m.Consume(ctx, systemID, amount)
}

func (m *Manager) systemIDByMessageID(messageID int) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, reminder := range m.reminders {
		if reminder.LastMessageID == messageID {
			return reminder.SystemID, true
		}
	}

	return "", false
}