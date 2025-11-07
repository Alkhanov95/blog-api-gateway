package logger

import (
	"fmt"
	"log/slog"
	"os"
)

const DefaultLevel = slog.LevelInfo

type Config struct {
	Level string `validate:"required,log_level"`
}

func Register(cfg Config) {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(cfg.Level)); err != nil {
		lvl = DefaultLevel
		slog.Warn(fmt.Sprintf("invalid min log level in cfg. Using default: %s",
			DefaultLevel.String()), slog.Any("error", err))
	} else {
		slog.Info("min log level set: " + lvl.String())
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))

	slog.SetDefault(logger)
}
