package telegram

import (
	"bytes"
	"context"
	"errors"
	"log"
	"testing"

	"umamusume-notifier/internal/app"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type mockService struct {
	setCalled         bool
	setSystemID       string
	setAmount         int
	setErr            error
	consumeReplyCalled bool
	consumeReplyMsgID  int
	consumeReplyAmount int
	consumeReplyErr    error
}

func (m *mockService) Status() []app.Status {
	return nil
}

func (m *mockService) Consume(context.Context, string, int) error {
	return nil
}

func (m *mockService) Set(
	_ context.Context,
	systemID string,
	amount int,
) error {
	m.setCalled = true
	m.setSystemID = systemID
	m.setAmount = amount
	return m.setErr
}

func (m *mockService) SetElapsed(context.Context, string, int) error {
	return nil
}

func (m *mockService) SetRegen(context.Context, string, int) error {
	return nil
}

func (m *mockService) ConsumeReply(_ context.Context, messageID int, amount int) error {
	m.consumeReplyCalled = true
	m.consumeReplyMsgID = messageID
	m.consumeReplyAmount = amount
	return m.consumeReplyErr
}

type mockSender struct {
	lastChatID int64
	lastText   string
}

func (m *mockSender) Send(chatID int64, text string) (int, error) {
	m.lastChatID = chatID
	m.lastText = text
	return 1, nil
}

func newTestBot(service Service) (*Bot, *mockSender) {
	sender := &mockSender{}
	var logBuf bytes.Buffer

	return &Bot{
		sender:  sender,
		service: service,
		logger:  log.New(&logBuf, "", 0),
	}, sender
}

func TestHandleSet(t *testing.T) {
	service := &mockService{}
	bot, sender := newTestBot(service)

	msg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 123},
		Text: "/set TP 20",
		Entities: []tgbotapi.MessageEntity{
			{
				Type:   "bot_command",
				Offset: 0,
				Length: 4,
			},
		},
	}

	bot.handleSet(msg)

	if !service.setCalled {
		t.Fatal("Set was not called")
	}

	if service.setSystemID != "TP" {
		t.Fatalf("systemID = %q, want %q", service.setSystemID, "TP")
	}

	if service.setAmount != 20 {
		t.Fatalf("amount = %d, want %d", service.setAmount, 20)
	}

	if sender.lastChatID != 123 {
		t.Fatalf("chatID = %d, want 123", sender.lastChatID)
	}

	if sender.lastText != "Updated TP: set 20 point(s)." {
		t.Fatalf("response = %q, want set confirmation", sender.lastText)
	}
}

func TestHandleSetError(t *testing.T) {
	service := &mockService{setErr: errors.New("boom")}
	bot, sender := newTestBot(service)

	msg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 123},
		Text: "/set TP 20",
		Entities: []tgbotapi.MessageEntity{
			{
				Type:   "bot_command",
				Offset: 0,
				Length: 4,
			},
		},
	}

	bot.handleCommand(msg)

	if sender.lastText != "boom" {
		t.Fatalf("response = %q, want %q", sender.lastText, "boom")
	}
}

func TestHandleReply_ReminderResponse(t *testing.T) {
	service := &mockService{}
	bot, sender := newTestBot(service)

	msg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 123},
		Text: "60",
		ReplyToMessage: &tgbotapi.Message{
			MessageID: 456,
			Text:      "✅ TP is now full.\n\nReply with the amount of points you use.",
		},
	}

	bot.handleReply(msg)

	if !service.consumeReplyCalled {
		t.Fatal("ConsumeReply was not called")
	}

	if service.consumeReplyMsgID != 456 {
		t.Fatalf("messageID = %d, want 456", service.consumeReplyMsgID)
	}

	if service.consumeReplyAmount != 60 {
		t.Fatalf("amount = %d, want 60", service.consumeReplyAmount)
	}

	if sender.lastText != "Recorded 60 point(s)." {
		t.Fatalf("response = %q, want %q", sender.lastText, "Recorded 60 point(s).")
	}
}

func TestHandleReply_InvalidAmount(t *testing.T) {
	service := &mockService{}
	bot, sender := newTestBot(service)

	msg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 123},
		Text: "sixty",
		ReplyToMessage: &tgbotapi.Message{
			MessageID: 456,
			Text:      "✅ TP is now full.\n\nReply with the amount of points you use.",
		},
	}

	bot.handleReply(msg)

	if service.consumeReplyCalled {
		t.Fatal("ConsumeReply should not be called")
	}

	if sender.lastText != "Reply with the number of points you used (for example: 20)." {
		t.Fatalf("response = %q, want validation message", sender.lastText)
	}
}
