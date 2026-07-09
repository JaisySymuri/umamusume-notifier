package telegram

import (
	"fmt"
	"log"

	// "umamusume-notifier/internal/notification"

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

