package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/filippofinke/docker-events/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, os.Stdout); err != nil {
		os.Exit(1)
	}
}
