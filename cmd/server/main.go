package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"umamusume-notifier/internal/app"
	"umamusume-notifier/internal/config"
	"umamusume-notifier/internal/metrics"
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

	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return

			case <-ticker.C:
				if _, err := os.Stat("shutdown.flag"); err == nil {
					log.Println("Shutdown flag detected.")
					stop()
					return
				} else if !os.IsNotExist(err) {
					log.Printf("Shutdown flag check failed: %v", err)
				}
			}
		}
	}()

	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	store, err := storage.NewSQLiteStore("data.db")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := store.Close(); err != nil {
			log.Printf("SQLite shutdown failed: %v", err)
		} else {
			log.Println("SQLite database connection closed.")
		}
	}()

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

	metricsServer := &http.Server{
		Addr:    metrics.ListenAddr(),
		Handler: metrics.Handler(),
	}

	go func() {
		log.Printf("metrics server listening on http://%s/metrics", metrics.ListenAddr())

		if err := metricsServer.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Printf("metrics server failed: %v", err)
			stop()
		}
	}()

	bot.SendText(cfg.Telegram.ChatID, telegram.FormatServiceOnline())

	go func() {
		log.Printf("metrics server listening on http://%s/metrics", metrics.ListenAddr())

		if err := metricsServer.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Printf("metrics server failed: %v", err)
			stop()
		}
	}()

	bot.SendText(cfg.Telegram.ChatID, telegram.FormatServiceOnline())

	if err := bot.Start(ctx); err != nil && ctx.Err() == nil {
		log.Printf("Bot failed: %v", err)
	}

	stop()

	log.Println("Shutting down application...")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := metricsServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Metrics server shutdown failed: %v", err)
	}

	log.Println("Application shutdown completed.")
}
