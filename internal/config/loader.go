package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	defaultSubjectPrefix = "Docker event"
)

func Load() (*Config, error) {
	cfg := &Config{
		DockerFilters:   splitAndTrim(os.Getenv("DOCKER_EVENT_FILTERS")),
		DockerEventType: splitAndTrim(os.Getenv("DOCKER_EVENT_TYPES")),
		NotifySubject:   getEnvOrDefault("NOTIFY_SUBJECT_PREFIX", defaultSubjectPrefix),
	}

	slackToken, ok := os.LookupEnv("SLACK_BOT_TOKEN")
	if ok && slackToken != "" {
		slackChannels := splitAndTrim(os.Getenv("SLACK_CHANNEL_IDS"))
		if len(slackChannels) == 0 {
			return nil, errors.New("slack configured but SLACK_CHANNEL_IDS is empty")
		}

		cfg.Slack = SlackConfig{
			Enabled:  true,
			Token:    slackToken,
			Channels: slackChannels,
		}
	}

	telegramToken, ok := os.LookupEnv("TELEGRAM_BOT_TOKEN")
	if ok && telegramToken != "" {
		rawChatIDs := splitAndTrim(os.Getenv("TELEGRAM_CHAT_IDS"))
		if len(rawChatIDs) == 0 {
			return nil, errors.New("telegram configured but TELEGRAM_CHAT_IDS is empty")
		}

		chatIDs := make([]int64, 0, len(rawChatIDs))
		for _, rawID := range rawChatIDs {
			chatID, err := strconv.ParseInt(rawID, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid TELEGRAM_CHAT_IDS value %q: %w", rawID, err)
			}
			chatIDs = append(chatIDs, chatID)
		}

		cfg.Telegram = TelegramConfig{
			Enabled: true,
			Token:   telegramToken,
			ChatIDs: chatIDs,
		}
	}

	discordToken, ok := os.LookupEnv("DISCORD_BOT_TOKEN")
	if ok && discordToken != "" {
		discordChannels := splitAndTrim(os.Getenv("DISCORD_CHANNEL_IDS"))
		if len(discordChannels) == 0 {
			return nil, errors.New("discord bot configured but DISCORD_CHANNEL_IDS is empty")
		}

		cfg.Discord = DiscordConfig{
			Enabled:    true,
			Token:      discordToken,
			ChannelIDs: discordChannels,
		}
	}

	discordWebhooks := splitAndTrim(os.Getenv("DISCORD_WEBHOOK_URLS"))
	if len(discordWebhooks) > 0 {
		if cfg.Discord.Enabled {
			cfg.Discord.WebhookURLs = discordWebhooks
		} else {
			cfg.Discord = DiscordConfig{
				Enabled:     true,
				WebhookURLs: discordWebhooks,
			}
		}
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getEnvOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}

	return fallback
}

func splitAndTrim(raw string) []string {
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
