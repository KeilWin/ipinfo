package app

import (
	"flag"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func setBootstrapLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)
}

func loadEnv() error {
	var envFile string
	flag.StringVar(&envFile, "env", dotEnvFilename(), "Environment filepath")
	flag.StringVar(&envFile, "e", dotEnvFilename(), "Environment filepath")
	flag.Parse()

	return godotenv.Load(envFile)
}

func bootstrap() (*IpInfoAppConfig, error) {
	setBootstrapLogger()
	if err := loadEnv(); err != nil {
		return nil, err
	}
	appConfig := NewIpInfoAppConfig()
	if err := appConfig.Load(); err != nil {
		return nil, err
	}
	if err := appConfig.Check(); err != nil {
		return nil, err
	}
	return appConfig, nil
}
