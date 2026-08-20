package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"umamusume-notifier/internal/metrics"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleCommand(msg *tgbotapi.Message) {
	switch msg.Command() {
	case "help":
		b.handleHelp(msg)

	case "status":
		fallthrough
	case "s":
		b.handleStatus(msg)
		metrics.ObserveCommand("status", "success")

	case "t":
		b.handleSoonestToFull(msg)
		metrics.ObserveCommand("t", "success")

	case "use":
		if err := b.handleUse(msg); err != nil {
			metrics.ObserveCommand("use", "error")
		} else {
			metrics.ObserveCommand("use", "success")
		}

	case "elapsed":
		if err := b.handleElapsed(msg); err != nil {
			metrics.ObserveCommand("elapsed", "error")
		} else {
			metrics.ObserveCommand("elapsed", "success")
		}

	case "regen":
		if err := b.handleRegen(msg); err != nil {
			metrics.ObserveCommand("regen", "error")
		} else {
			metrics.ObserveCommand("regen", "success")
		}
	case "set":
		if err := b.handleSet(msg); err != nil {
			metrics.ObserveCommand("set", "error")
		} else {
			metrics.ObserveCommand("set", "success")
		}

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

func (b *Bot) handleSoonestToFull(msg *tgbotapi.Message) {
	status := b.service.Status()
	b.SendText(msg.Chat.ID, FormatSoonestToFull(status))
}

func (b *Bot) handleUse(msg *tgbotapi.Message) error {
	systemID, amount, err := ParseUse(msg.CommandArguments())
	if err != nil {
		b.SendText(msg.Chat.ID, err.Error())
		return err
	}

	if err := b.service.Consume(
		context.Background(),
		systemID,
		amount,
	); err != nil {
		b.SendText(msg.Chat.ID, err.Error())
		return err
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

	return nil
}

func (b *Bot) handleSet(msg *tgbotapi.Message) error {
	systemID, amount, err := ParseSet(msg.CommandArguments())
	if err != nil {
		b.SendText(msg.Chat.ID, err.Error())
		return err
	}

	if err := b.service.Set(
		context.Background(),
		systemID,
		amount,
	); err != nil {
		b.SendText(msg.Chat.ID, err.Error())
		return err
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

	return nil
}

func (b *Bot) handleElapsed(msg *tgbotapi.Message) error {
	systemID, minutes, err := ParseElapsed(msg.CommandArguments())
	if err != nil {
		b.SendText(msg.Chat.ID, err.Error())
		return err
	}

	if err := b.service.SetElapsed(
		context.Background(),
		systemID,
		minutes,
	); err != nil {
		b.SendText(msg.Chat.ID, err.Error())
		return err
	}

	b.SendText(
		msg.Chat.ID,
		fmt.Sprintf(
			"Updated %s: elapsed time set to %d minute(s).",
			systemID,
			minutes,
		),
	)

	return nil
}

func (b *Bot) handleRegen(msg *tgbotapi.Message) error {
	systemID, minutesLeft, err := ParseRegen(msg.CommandArguments())
	if err != nil {
		b.SendText(msg.Chat.ID, err.Error())
		return err
	}

	if err := b.service.SetRegen(
		context.Background(),
		systemID,
		minutesLeft,
	); err != nil {
		b.SendText(msg.Chat.ID, err.Error())
		return err
	}

	b.SendText(
		msg.Chat.ID,
		fmt.Sprintf(
			"Updated %s: %d minute(s) left until the next point.",
			systemID,
			minutesLeft,
		),
	)

	return nil
}

func (b *Bot) handleReply(msg *tgbotapi.Message) {
	amount, err := strconv.Atoi(strings.TrimSpace(msg.Text))
	if err != nil {
		b.SendText(
			msg.Chat.ID,
			"Reply with the number of points you used (for example: 20).",
		)
		metrics.ObserveCommand("reply", "error")
		return
	}

	if err := b.service.ConsumeReply(
		context.Background(),
		msg.ReplyToMessage.MessageID,
		amount,
	); err != nil {
		b.SendText(msg.Chat.ID, err.Error())
		metrics.ObserveCommand("reply", "error")
		return
	}

	b.SendText(
		msg.Chat.ID,
		fmt.Sprintf("Recorded %d point(s).", amount),
	)
	metrics.ObserveCommand("reply", "success")
}

func (b *Bot) handleUnknownCommand(msg *tgbotapi.Message) {
	b.SendText(msg.Chat.ID, "Unknown command. Type /help for a list of available commands.")
}
