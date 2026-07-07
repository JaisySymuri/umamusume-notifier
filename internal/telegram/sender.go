package telegram

import (
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

