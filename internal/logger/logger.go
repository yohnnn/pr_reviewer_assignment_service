package logger

import (
	"log/slog"
	"os"
)

func New(env, level string) *slog.Logger {
	var lvl slog.Level

	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     lvl,
		AddSource: true,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)

	return slog.New(handler)
}
