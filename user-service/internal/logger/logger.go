package logger

import (
	"log/slog"
	"os"
	"sync"
)

var (
	Log  *slog.Logger
	once sync.Once
)

// Init инициализирует логгер в зависимости от окружения: "DEV" (текст) или "PROD" (JSON)
func Init(env string) {
	once.Do(func() {
		var handler slog.Handler

		switch env {
		case "DEV":
			handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			})
		default: // "PROD" или любое другое
			handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			})
		}

		Log = slog.New(handler)
	})
}
