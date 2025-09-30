package logging

import (
	"io"
	"log/slog"
)

func NewLogger(out io.Writer) *slog.Logger {
	return slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{Level: slog.LevelInfo}))
}
