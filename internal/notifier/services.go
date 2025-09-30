package notifier

import (
	"fmt"
	"strings"

	"github.com/nikoksr/notify/service/discord"
	"github.com/nikoksr/notify/service/slack"
	"github.com/nikoksr/notify/service/telegram"

	"github.com/filippofinke/docker-events/internal/config"
)

func (n *notifierImpl) addSlack(cfg config.SlackConfig) error {
	if strings.TrimSpace(cfg.Token) == "" {
		return fmt.Errorf("empty slack token")
	}
	if len(cfg.Channels) == 0 {
		return fmt.Errorf("no slack channels configured")
	}
	service := slack.New(cfg.Token)
	service.AddReceivers(cfg.Channels...)
	n.client.UseServices(service)
	return nil
}

func (n *notifierImpl) addTelegram(cfg config.TelegramConfig) error {
	service, err := telegram.New(cfg.Token)
	if err != nil {
		return fmt.Errorf("create telegram service: %w", err)
	}
	service.SetParseMode("")
	service.AddReceivers(cfg.ChatIDs...)
	n.client.UseServices(service)
	return nil
}

func (n *notifierImpl) addDiscord(cfg config.DiscordConfig) error {
	if strings.TrimSpace(cfg.Token) == "" {
		return fmt.Errorf("empty discord token")
	}
	if len(cfg.ChannelIDs) == 0 {
		return fmt.Errorf("no discord channels configured")
	}
	service := discord.New()
	if err := service.AuthenticateWithBotToken(cfg.Token); err != nil {
		return fmt.Errorf("authenticate discord bot: %w", err)
	}
	service.AddReceivers(cfg.ChannelIDs...)
	n.client.UseServices(service)
	return nil
}
