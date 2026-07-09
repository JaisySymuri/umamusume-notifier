package telegram

import (
	"umamusume-notifier/internal/notification"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Sender interface {
	Send(chatID int64, text string) (int, error)
}

type telegramSender struct {
	api *tgbotapi.BotAPI
}

func (s *telegramSender) Send(
	chatID int64,
	text string,
) (int, error) {
	msg := tgbotapi.NewMessage(chatID, text)

	sent, err := s.api.Send(msg)
	if err != nil {
		return 0, err
	}

	return sent.MessageID, nil
}

func (b *Bot) SendText(chatID int64, text string) {
	_, err := b.sender.Send(chatID, text)
	if err != nil {
		b.logger.Printf("telegram send failed: chat_id=%d: %v", chatID, err)
	}
}

func (b *Bot) SendNotification(chatID int64, event notification.Event) (int, error) {
	messageID, err := b.sender.Send(chatID, FormatNotification(event))
	if err != nil {
		b.logger.Printf("telegram notify failed: chat_id=%d: %v", chatID, err)
		return 0, err
	}

	return messageID, nil
}

// Notify adapts Bot to the notification.Sender interface.
func (b *Bot) Notify(chatID int64, event notification.Event) (int, error) {
	return b.SendNotification(chatID, event)
}
