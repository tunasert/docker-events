package app

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/filippofinke/docker-events/internal/config"
	"github.com/filippofinke/docker-events/internal/docker"
	"github.com/filippofinke/docker-events/internal/logging"
	"github.com/filippofinke/docker-events/internal/notifier"
)

func Run(ctx context.Context, logOut io.Writer) error {
	logger := logging.NewLogger(logOut)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("load configuration", "error", err)
		return fmt.Errorf("load config: %w", err)
	}

	n := notifier.NewNotifier(logger)
	if err := n.Setup(cfg); err != nil {
		logger.Error("configure notifier", "error", err)
		return fmt.Errorf("setup notifier: %w", err)
	}

	watcher, err := docker.NewWatcher(cfg.DockerFilters, cfg.DockerEventType, logger)
	if err != nil {
		logger.Error("create docker watcher", "error", err)
		return fmt.Errorf("create watcher: %w", err)
	}
	logger.Info("starting docker events watcher", "filters", cfg.DockerFilters, "types", cfg.DockerEventType)

	err = watcher.Watch(ctx, func(ctx context.Context, event docker.Event) error {
		attrs := make([]any, 0, 10)
		attrs = append(attrs, "type", event.Type, "action", event.Action, "status", event.Status, "id", event.ID)
		if event.Actor.ID != "" {
			attrs = append(attrs, "actor", event.Actor.ID)
		}
		attrs = append(attrs, "timestamp", event.Timestamp.Format("2006-01-02T15:04:05Z07:00"))
		logger.Info("docker event", attrs...)
		return n.NotifyEvent(ctx, cfg, event)
	})

	if err != nil {
		if errors.Is(err, context.Canceled) {
			logger.Info("docker events watcher stopped", "reason", "context cancelled")
			return nil
		}
		logger.Error("watch docker events", "error", err)
		return err
	}

	return nil
}
