package notifier

import (
	"context"
	"fmt"

	"github.com/nikoksr/notify"

	"github.com/filippofinke/docker-events/internal/config"
	"github.com/filippofinke/docker-events/internal/docker"

	"log/slog"
)

type Notifier interface {
	Setup(cfg *config.Config) error
	NotifyEvent(ctx context.Context, cfg *config.Config, event docker.Event) error
}

type notifierImpl struct {
	client *notify.Notify
	logger *slog.Logger
}

func NewNotifier(logger *slog.Logger) Notifier {
	return &notifierImpl{
		client: notify.New(),
		logger: logger,
	}
}

func (n *notifierImpl) Setup(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("nil config")
	}

	if cfg.Slack.Enabled {
		if err := n.addSlack(cfg.Slack); err != nil {
			return fmt.Errorf("setup slack: %w", err)
		}
	}

	if cfg.Telegram.Enabled {
		if err := n.addTelegram(cfg.Telegram); err != nil {
			return fmt.Errorf("setup telegram: %w", err)
		}
	}

	if cfg.Discord.Enabled {
		if err := n.addDiscord(cfg.Discord); err != nil {
			return fmt.Errorf("setup discord: %w", err)
		}
	}

	return nil
}

func (n *notifierImpl) NotifyEvent(ctx context.Context, cfg *config.Config, event docker.Event) error {
	if cfg == nil {
		return fmt.Errorf("nil config")
	}
	subject, body := formatEvent(cfg.NotifySubject, event)
	if err := n.client.Send(ctx, subject, body); err != nil {
		return fmt.Errorf("send notification: %w", err)
	}
	return nil
}
