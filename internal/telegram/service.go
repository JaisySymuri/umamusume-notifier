package telegram

import (
	"context"

	"umamusume-notifier/internal/app"
)

type Service interface {
	Status() []app.Status

	Consume(
		ctx context.Context,
		systemID string,
		amount int,
	) error

	SetElapsed(
		ctx context.Context,
		systemID string,
		minutes int,
	) error

	ConsumeReply(
		ctx context.Context,
		messageID int,
		amount int,
	) error
}
