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
	setCalled   bool
	setSystemID string
	setAmount   int
	setErr      error
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

func (m *mockService) ConsumeReply(context.Context, int, int) error {
	return nil
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
