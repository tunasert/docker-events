package config

import (
	"strings"
	"time"
)

type Config struct {
	DockerFilters    []string
	DockerEventType  []string
	NotifySubject    string
	MessageTemplate  string
	LogLines         int
	EventGroupWindow time.Duration
	Slack            SlackConfig
	Telegram         TelegramConfig
	Discord          DiscordConfig
}

type SlackConfig struct {
	Enabled  bool
	Token    string
	Channels []string
}

type TelegramConfig struct {
	Enabled bool
	Token   string
	ChatIDs []int64
}

type DiscordConfig struct {
	Enabled     bool
	Token       string
	ChannelIDs  []string
	WebhookURLs []string
}

func (c *Config) Validate() error {
	var missing []string

	if !c.Slack.Enabled && !c.Telegram.Enabled && !c.Discord.Enabled {
		missing = append(missing, "notification credentials (Slack, Telegram, or Discord)")
	}

	if len(missing) > 0 {
		return &configError{msg: "missing required configuration: " + strings.Join(missing, ", ")}
	}

	return nil
}

type configError struct{ msg string }

func (e *configError) Error() string { return e.msg }
