package telegram

import (
	"fmt"
	"log"

	"umamusume-notifier/internal/notification"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot wraps the Telegram API client used by the application.
type Bot struct {
	api     *tgbotapi.BotAPI
	sender  Sender
	service Service
	logger  *log.Logger
}

// New creates and validates a Telegram bot client.
func New(
	token string,
	service Service,
	logger *log.Logger,
) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("create telegram bot: %w", err)
	}

	return &Bot{
		api:     api,
		sender:  &telegramSender{api: api},
		service: service,
		logger:  logger,
	}, nil
}

func (b *Bot) send(chatID int64, text string) {
	_, err := b.sender.Send(chatID, text)
	if err != nil {
		b.logger.Printf("telegram send failed: chat_id=%d: %v", chatID, err)
	}
}

func (b *Bot) Notify(chatID int64, event notification.Event) (int, error) {
	messageID, err := b.sender.Send(chatID, FormatNotification(event))
	if err != nil {
		b.logger.Printf("telegram send failed: chat_id=%d: %v", chatID, err)
		return 0, err
	}

	return messageID, nil
}
