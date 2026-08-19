package logger

import (
	"log/slog"
	"os"
)

// New создает новый логгер с указанным уровнем
func InitLogger(confLevel string) *slog.Logger {
	// Устанавливаем уровень по умолчанию
	level := slog.LevelInfo

	switch confLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		// Если уровень не распознан, используем INFO
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	return logger
}
