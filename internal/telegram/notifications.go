package telegram

import (
	"context"
	"fmt"

	"umamusume-notifier/internal/notification"
)

type Notifier struct {
	bot     *Bot
	manager interface {
		RecordReminderMessage(
			context.Context,
			string,
			int,
		) error
	}

	chatID int64
}

// Notify sends a reminder notification and returns the Telegram message ID.
func (n *Notifier) Notify(
	ctx context.Context,
	event notification.Event,
) error {
	messageID, err := n.bot.SendNotification(
		n.chatID,
		event,
	)
	if err != nil {
		return err
	}

	return n.manager.RecordReminderMessage(
		ctx,
		event.SystemID,
		messageID,
	)
}

// FormatNotification formats a reminder notification.
func FormatNotification(event notification.Event) string {
	switch event.Type {
	case notification.NearFull:
		return fmt.Sprintf(
			"⚠️ %s is almost full.\n\nReply with the amount of points you use.",
			event.SystemID,
		)

	case notification.Full:
		return fmt.Sprintf(
			"✅ %s is now full.\n\nReply with the amount of points you use.",
			event.SystemID,
		)

	default:
		return event.SystemID
	}
}
