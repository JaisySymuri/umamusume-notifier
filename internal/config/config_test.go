package config

import "testing"

func TestConfigValidate_DuplicateSystemID(t *testing.T) {
	cfg := Config{
		Telegram: TelegramConfig{
			Token:  "token",
			ChatID: 123,
		},
		Scheduler: SchedulerConfig{
			TickInterval:  1,
			AlertThreshold: 1,
		},
		Systems: []PointSystemConfig{
			{
				ID:           "TP",
				Name:         "Training Points",
				Max:          100,
				RegenMinutes: 10,
			},
			{
				ID:           "TP",
				Name:         "Duplicate Training Points",
				Max:          50,
				RegenMinutes: 20,
			},
		},
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected duplicate system ID validation error")
	}
}
