package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Telegram TelegramConfig `yaml:"telegram"`
	Scheduler SchedulerConfig `yaml:"scheduler"`
	Systems []PointSystemConfig `yaml:"systems"`
}

type TelegramConfig struct {
	Token  string `yaml:"token"`
	ChatID int64  `yaml:"chat_id"`
}

type SchedulerConfig struct {
	TickInterval  time.Duration `yaml:"tick_interval"`
	AlertThreshold time.Duration `yaml:"alert_threshold"`
}

type PointSystemConfig struct {
	ID            string `yaml:"id"`
	Name          string `yaml:"name"`
	Max           int    `yaml:"max"`
	RegenMinutes  int    `yaml:"regen_minutes"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.Telegram.Token == "" {
		return fmt.Errorf("telegram.token is required")
	}

	if c.Telegram.ChatID == 0 {
		return fmt.Errorf("telegram.chat_id is required")
	}

	if c.Scheduler.TickInterval <= 0 {
		return fmt.Errorf("scheduler.tick_interval must be greater than zero")
	}

	if c.Scheduler.AlertThreshold <= 0 {
		return fmt.Errorf("scheduler.alert_threshold must be greater than zero")
	}

	if len(c.Systems) == 0 {
		return fmt.Errorf("at least one point system must be configured")
	}

	for _, system := range c.Systems {
		if system.ID == "" {
			return fmt.Errorf("system.id is required")
		}

		if system.Name == "" {
			return fmt.Errorf("system.name is required")
		}

		if system.Max <= 0 {
			return fmt.Errorf("%s: max must be greater than zero", system.ID)
		}

		if system.RegenMinutes <= 0 {
			return fmt.Errorf("%s: regen_minutes must be greater than zero", system.ID)
		}
	}

	return nil
}