package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"umamusume-notifier/internal/app"
	gas "umamusume-notifier/internal/appdynamics"
	"umamusume-notifier/internal/config"
	"umamusume-notifier/internal/metrics"
	"umamusume-notifier/internal/notification"
	"umamusume-notifier/internal/points"
	"umamusume-notifier/internal/scheduler"
	"umamusume-notifier/internal/storage"
	"umamusume-notifier/internal/telegram"
)

func main() {
	acfg := gas.Config{}
	// Controller
	acfg.Controller.Host = "192.168.1.124.nip.io"
	acfg.Controller.Port = 443
	acfg.Controller.UseSSL = true
	acfg.Controller.Account = "customer1"
	acfg.Controller.AccessKey = "d6455146-9662-4e13-bae3-fe7958fd1ea6"
	// App Context
	acfg.AppName = "telegram-bot-app"
	acfg.TierName = "telegram-bot-tier"
	acfg.NodeName = "telegram-bot-node"
	// misc
	acfg.InitTimeoutMs = 1000
	// init the SDK - Only for Linux
	if err := gas.InitSDK(&acfg); err != nil {
		fmt.Printf("Error initializing the AppDynamics SDK\n")
	} else {
		fmt.Printf("Initialized AppDynamics SDK successfully\n")
	}

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

	// creates an empty slice of points.Definition with enough capacity to hold all systems from config
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
		log.Default(),
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
	go func() {
		log.Printf("metrics server listening on http://%s/metrics", metrics.ListenAddr())
		if err := http.ListenAndServe(metrics.ListenAddr(), metrics.Handler()); err != nil && ctx.Err() == nil {
			log.Printf("metrics server failed: %v", err)
		}
	}()

	bot.SendText(cfg.Telegram.ChatID, telegram.FormatServiceOnline())

	if err := bot.Start(ctx); err != nil && ctx.Err() == nil {
		log.Fatal(err)
	}

	// Stop/Clean up the AppD SDK. Bottom of the func
	gas.TerminateSDK()

}
