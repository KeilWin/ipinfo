package app

import (
	"log/slog"
	"os"
)

func NewAppLoggerHandler(appCfg *IpInfoAppConfig) *slog.JSONHandler {
	return slog.NewJSONHandler(os.Stdout, nil)
}

func NewAppLogger(appCfg *IpInfoAppConfig) *slog.Logger {
	logger := slog.New(NewAppLoggerHandler(appCfg))
	slog.SetDefault(logger)
	return logger
}
