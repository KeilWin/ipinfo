package logger

import (
	"log/slog"
	"os"

	"github.com/KeilWin/ipinfo/internal/common"
)

type LoggerConfig struct {
	common.Config
}

func (p *LoggerConfig) Load() error {
	return nil
}

func (p *LoggerConfig) Check() error {
	return nil
}

func NewLoggerConfig() *LoggerConfig {
	return &LoggerConfig{}
}

func NewAppLoggerHandler(config *LoggerConfig) *slog.JSONHandler {
	return slog.NewJSONHandler(os.Stdout, nil)
}

func NewAppLogger(config *LoggerConfig) *slog.Logger {
	logger := slog.New(NewAppLoggerHandler(config))
	slog.SetDefault(logger)
	return logger
}
