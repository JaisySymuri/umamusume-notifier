package notification

import (
	"context"
)

type Sender interface {
	Notify(
		chatID int64,
		event Event,
	) (int, error)
}

type Recorder interface {
	RecordReminderMessage(
		ctx context.Context,
		systemID string,
		messageID int,
	) error
}

type Service struct {
	sender   Sender
	recorder Recorder
	chatID   int64
}

func NewService(
	sender Sender,
	recorder Recorder,
	chatID int64,
) *Service {
	return &Service{
		sender:   sender,
		recorder: recorder,
		chatID:   chatID,
	}
}

func (s *Service) Notify(
	ctx context.Context,
	event Event,
) error {
	messageID, err := s.sender.Notify(
		s.chatID,
		event,
	)
	if err != nil {
		return err
	}

	return s.recorder.RecordReminderMessage(
		ctx,
		event.SystemID,
		messageID,
	)
}