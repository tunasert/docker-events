package notifier

import (
	"context"
	"fmt"

	"github.com/nikoksr/notify"

	"github.com/filippofinke/docker-events/internal/config"
	"github.com/filippofinke/docker-events/internal/docker"

	"log/slog"

	dockerclient "github.com/docker/docker/client"
)

type Notifier interface {
	Setup(cfg *config.Config) error
	NotifyEvent(ctx context.Context, cfg *config.Config, event docker.Event) error
	SetDockerClient(cli *dockerclient.Client)
}

type notifierImpl struct {
	client    *notify.Notify
	logger    *slog.Logger
	dockerCli *dockerclient.Client
}

func NewNotifier(logger *slog.Logger) Notifier {
	return &notifierImpl{
		client: notify.New(),
		logger: logger,
	}
}

func (n *notifierImpl) SetDockerClient(cli *dockerclient.Client) {
	n.dockerCli = cli
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

	var subject, body string
	var err error

	if cfg.MessageTemplate != "" {
		body, _, err = formatEventWithTemplate(cfg.MessageTemplate, event, n.dockerCli, cfg.LogLines)
		if err != nil {
			n.logger.Warn("failed to format event with template, falling back to default format", "error", err)
			subject, body = formatEvent(cfg.NotifySubject, event)
		} else {
			subject = fmt.Sprintf("%s: %s %s", cfg.NotifySubject, event.Type, event.Action)
		}
	} else {
		subject, body = formatEvent(cfg.NotifySubject, event)
	}

	if err := n.client.Send(ctx, subject, body); err != nil {
		return fmt.Errorf("send notification: %w", err)
	}
	return nil
}
