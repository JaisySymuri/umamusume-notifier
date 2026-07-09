package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"umamusume-notifier/internal/app"
	"umamusume-notifier/internal/config"
	"umamusume-notifier/internal/notification"
	"umamusume-notifier/internal/points"
	"umamusume-notifier/internal/scheduler"
	"umamusume-notifier/internal/storage"
	"umamusume-notifier/internal/telegram"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	store, err := storage.NewSQLiteStore("data.db")
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Initialize(ctx); err != nil {
		log.Fatal(err)
	}

	definitions := make([]points.Definition, 0, len(cfg.Systems))

	for _, system := range cfg.Systems {
		definitions = append(definitions, points.Definition{
			ID:           system.ID,
			Name:         system.Name,
			Max:          system.Max,
			RegenMinutes: system.RegenMinutes,
		})
	}

	manager := app.New(
		store,
		cfg.Scheduler.AlertThreshold,
	)

	if err := manager.Load(ctx, definitions); err != nil {
		log.Fatal(err)
	}

	bot, err := telegram.New(
		cfg.Telegram.Token,
		manager,
		log.Default(),
	)
	if err != nil {
		log.Fatal(err)
	}

	notificationService := notification.NewService(
		bot,
		manager,
		cfg.Telegram.ChatID,
	)

	scheduler := scheduler.New(
		manager,
		notificationService,
		cfg.Scheduler.TickInterval,
		log.Default(),
	)

	go scheduler.Run(ctx)

	bot.SendText(cfg.Telegram.ChatID, telegram.FormatServiceOnline())

	if err := bot.Start(ctx); err != nil && ctx.Err() == nil {
		log.Fatal(err)
	}

}
