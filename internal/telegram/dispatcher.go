package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Start begins receiving updates from Telegram.
//
// It blocks until the context is cancelled.
func (b *Bot) Start(ctx context.Context) error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := b.api.GetUpdatesChan(updateConfig)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case update := <-updates:
			b.dispatch(update)
		}
	}
}

// dispatch routes an incoming update to the appropriate handler.
func (b *Bot) dispatch(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	if update.Message.ReplyToMessage != nil {
		b.handleReply(update.Message)
		return
	}

	if update.Message.IsCommand() {
		b.handleCommand(update.Message)
		return
	}
}