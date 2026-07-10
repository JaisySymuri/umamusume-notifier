package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleCommand(msg *tgbotapi.Message) {
	switch msg.Command() {
	case "help":
		b.handleHelp(msg)

	case "status":
		b.handleStatus(msg)

	case "use":
		b.handleUse(msg)

	case "elapsed":
		b.handleElapsed(msg)

	case "regen":
		b.handleRegen(msg)
	
	case "set":
		b.handleSet(msg)

	default:
		b.handleUnknownCommand(msg)
	}
}

func (b *Bot) handleHelp(msg *tgbotapi.Message) {
	b.SendText(msg.Chat.ID, FormatHelp())
}

func (b *Bot) handleStatus(msg *tgbotapi.Message) {
	status := b.service.Status()
	b.SendText(msg.Chat.ID, FormatStatus(status))
}

func (b *Bot) handleUse(msg *tgbotapi.Message) {
	systemID, amount, err := ParseUse(msg.CommandArguments())
	if err != nil {
		b.SendText(msg.Chat.ID, err.Error())
		return
	}

	if err := b.service.Consume(
		context.Background(),
		systemID,
		amount,
	); err != nil {
		b.SendText(msg.Chat.ID, err.Error())
		return
	}

	action := "consumed"
	if amount < 0 {
		action = "added"
		amount = -amount
	}

	b.SendText(
		msg.Chat.ID,
		fmt.Sprintf(
			"Updated %s: %s %d point(s).",
			systemID,
			action,
			amount,
		),
	)
}

func (b *Bot) handleSet(msg *tgbotapi.Message) {
	systemID, amount, err := ParseSet(msg.CommandArguments())
	if err != nil {
		b.SendText(msg.Chat.ID, err.Error())
		return
	}

	if err := b.service.Set(
		context.Background(),
		systemID,
		amount,
	); err != nil {
		b.SendText(msg.Chat.ID, err.Error())
		return
	}

	action := "set"
	if amount < 0 {
		action = "added"
		amount = -amount
	}

	b.SendText(
		msg.Chat.ID,
		fmt.Sprintf(
			"Updated %s: %s %d point(s).",
			systemID,
			action,
			amount,
		),
	)
}

func (b *Bot) handleElapsed(msg *tgbotapi.Message) {
	systemID, minutes, err := ParseElapsed(msg.CommandArguments())
	if err != nil {
		b.SendText(msg.Chat.ID, err.Error())
		return
	}

	if err := b.service.SetElapsed(
		context.Background(),
		systemID,
		minutes,
	); err != nil {
		b.SendText(msg.Chat.ID, err.Error())
		return
	}

	b.SendText(
		msg.Chat.ID,
		fmt.Sprintf(
			"Updated %s: elapsed time set to %d minute(s).",
			systemID,
			minutes,
		),
	)
}

func (b *Bot) handleRegen(msg *tgbotapi.Message) {
	systemID, minutesLeft, err := ParseRegen(msg.CommandArguments())
	if err != nil {
		b.SendText(msg.Chat.ID, err.Error())
		return
	}

	if err := b.service.SetRegen(
		context.Background(),
		systemID,
		minutesLeft,
	); err != nil {
		b.SendText(msg.Chat.ID, err.Error())
		return
	}

	b.SendText(
		msg.Chat.ID,
		fmt.Sprintf(
			"Updated %s: %d minute(s) left until the next point.",
			systemID,
			minutesLeft,
		),
	)
}



func (b *Bot) handleReply(msg *tgbotapi.Message) {
	amount, err := strconv.Atoi(strings.TrimSpace(msg.Text))
	if err != nil {
		b.SendText(
			msg.Chat.ID,
			"Reply with the number of points you used (for example: 20).",
		)
		return
	}

	if err := b.service.ConsumeReply(
		context.Background(),
		msg.ReplyToMessage.MessageID,
		amount,
	); err != nil {
		b.SendText(msg.Chat.ID, err.Error())
		return
	}

	b.SendText(
		msg.Chat.ID,
		fmt.Sprintf("Recorded %d point(s).", amount),
	)
}

func (b *Bot) handleUnknownCommand(msg *tgbotapi.Message) {
	b.SendText(msg.Chat.ID, "Unknown command. Type /help for a list of available commands.")
}
