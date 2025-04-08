package slogger

import (
	"log/slog"
	"os"
	"strings"
)

var logger *slog.Logger

func init() {
	logLevel := parseLogLevel(os.Getenv("LOG_LEVEL"))

	handler := slog.NewJSONHandler(
		os.Stdout, &slog.HandlerOptions{Level: logLevel},
	)

	logger = slog.New(handler)
}

func PkgLogger(pkg string) *slog.Logger {
	return logger.With(
		slog.String("package", pkg),
	)
}

func parseLogLevel(level string) slog.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
