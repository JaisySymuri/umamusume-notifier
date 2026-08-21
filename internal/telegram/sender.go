package telegram

import (
	"context"
	"errors"
	"net"
	"strings"
	"time"

	"umamusume-notifier/internal/metrics"
	"umamusume-notifier/internal/notification"
	gas "umamusume-notifier/internal/appdynamics"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Sender interface {
	Send(chatID int64, text string) (int, error)
}

type telegramAPI interface {
	Send(tgbotapi.Chattable) (tgbotapi.Message, error)
	GetUpdatesChan(tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
}

type telegramSender struct {
	api telegramAPI
}

func (s *telegramSender) Send(
	chatID int64,
	text string,
) (int, error) {
	start := time.Now()
	msg := tgbotapi.NewMessage(chatID, text)

	sent, err := s.api.Send(msg)
	if err != nil {
		metrics.ObserveTelegramAPIRequest("sendMessage", "error", time.Since(start))
		metrics.ObserveTelegramAPIError("sendMessage", classifyTelegramError(err))
		return 0, err
	}

	metrics.ObserveTelegramAPIRequest("sendMessage", "success", time.Since(start))

	return sent.MessageID, nil
}

func (b *Bot) SendText(chatID int64, text string) {
	btHandle := gas.StartBT("SendText", "")

    defer func() {
        gas.EndBT(btHandle)
    }()

	_, err := b.sender.Send(chatID, text)
	if err != nil {
		gas.AddBTError(
                btHandle,
                gas.APPD_LEVEL_ERROR,
                err.Error(),
                true,
            )
		b.logger.Printf("telegram send failed: chat_id=%d: %v", chatID, err)
	}
}

func (b *Bot) SendNotification(chatID int64, event notification.Event) (int, error) {
	// start the "Checkout" transaction. On top of func
    
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

func classifyTelegramError(err error) string {
	if err == nil {
		return "unknown"
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		if netErr.Timeout() {
			return "timeout"
		}
		return "network"
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return "timeout"
	}

	msg := strings.ToLower(err.Error())

	switch {
	case strings.Contains(msg, "too many requests"), strings.Contains(msg, "429"):
		return "rate_limit"
	case strings.Contains(msg, "telegram"), strings.Contains(msg, "api response"):
		return "telegram_api"
	case strings.Contains(msg, "timeout"):
		return "timeout"
	case strings.Contains(msg, "connection"), strings.Contains(msg, "network"), strings.Contains(msg, "no such host"), strings.Contains(msg, "temporary"):
		return "network"
	default:
		return "unknown"
	}
}
