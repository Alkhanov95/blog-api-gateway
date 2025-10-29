package main

import (
	"log/slog"
	"os"
)

func main() {
	// Запуск приложения
	if err := app.Run(); err != nil {
		slog.Error("failed to start app", slog.Any("error", err))
		os.Exit(1)
	}
}
