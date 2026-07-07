package scheduler

import (
	"context"
	"log"
	"time"

	"umamusume-notifier/internal/notification"
)

type TickManager interface {
	Tick(ctx context.Context, now time.Time) ([]notification.Event, error)
}

type Scheduler struct {
	manager  TickManager
	notifier Notifier

	ticker *time.Ticker
	logger *log.Logger
}

type Notifier interface {
	Notify(
		ctx context.Context,
		event notification.Event,
	) error
}

func New(
	manager TickManager,
	notifier Notifier,
	interval time.Duration,
	logger *log.Logger,
) *Scheduler {

	return &Scheduler{
		manager:  manager,
		notifier: notifier,
		ticker:   time.NewTicker(interval),
		logger:   logger,
	}
}

func (s *Scheduler) Run(ctx context.Context) {

	defer s.ticker.Stop()

	for {
		select {

		case <-ctx.Done():
			return

		case now := <-s.ticker.C:

			events, err := s.manager.Tick(ctx, now)
			if err != nil {
				s.logger.Printf("scheduler tick failed: %v", err)
				continue
			}

			for _, event := range events {
				if err := s.notifier.Notify(ctx, event); err != nil {
					s.logger.Printf(
						"notify %s failed: %v",
						event.SystemID,
						err,
					)
				}
			}
		}
	}
}
