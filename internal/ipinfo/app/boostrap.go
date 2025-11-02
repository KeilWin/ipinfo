package app

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"github.com/KeilWin/ipinfo/internal/ipinfo/handler"
)

const AppName = "ipinfo"

type ProtocolType string

const (
	ProtocolHTTP  ProtocolType = "http"
	ProtocolHTTPS ProtocolType = "https"
)

type IpInfoAppConfig struct {
	Addr              string
	Protocol          ProtocolType
	MaxHeaderBytes    int
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration

	CertFile string
	KeyFile  string

	HandlerConfig handler.AppHandlerConfig
}

func dotEnvFilename() string {
	return fmt.Sprintf("./configs/%s.env", AppName)
}

func setBootstrapLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)
}

func loadConfig() *IpInfoAppConfig {
	var envFile string
	flag.StringVar(&envFile, "env", dotEnvFilename(), "Environment filepath")
	flag.StringVar(&envFile, "e", dotEnvFilename(), "Environment filepath")
	flag.Parse()

	err := godotenv.Load(envFile)
	if err != nil {
		slog.Error("can't load env file", "file", envFile, "error", err)
		os.Exit(-1)
	}

	addr := os.Getenv("IPINFO_ADDR")
	protocol := ProtocolType(os.Getenv("IPINFO_PROTOCOL"))
	maxHeaderBytes, err := strconv.Atoi(os.Getenv("IPINFO_MAX_HEADER_BYTES"))
	if err != nil {
		slog.Error("can't init config", "field", "IPINFO_MAX_HEADER_BYTES", "error", err)
		os.Exit(-1)
	}
	readTimeout, err := strconv.Atoi(os.Getenv("IPINFO_READ_TIMEOUT"))
	if err != nil {
		slog.Error("can't init config", "field", "IPINFO_READ_TIMEOUT", "error", err)
		os.Exit(-1)
	}
	readHeaderTimeout, err := strconv.Atoi(os.Getenv("IPINFO_READ_HEADER_TIMEOUT"))
	if err != nil {
		slog.Error("can't init config", "field", "IPINFO_READ_HEADER_TIMEOUT", "error", err)
		os.Exit(-1)
	}
	writeTimeout, err := strconv.Atoi(os.Getenv("IPINFO_WRITE_TIMEOUT"))
	if err != nil {
		slog.Error("can't init config", "field", "IPINFO_WRITE_TIMEOUT", "error", err)
		os.Exit(-1)
	}
	idleTimeout, err := strconv.Atoi(os.Getenv("IPINFO_IDLE_TIMEOUT"))
	if err != nil {
		slog.Error("can't init config", "field", "IPINFO_IDLE_TIMEOUT", "error", err)
		os.Exit(-1)
	}
	certFile := os.Getenv("IPINFO_CERT_FILE")
	keyFile := os.Getenv("IPINFO_KEY_FILE")

	baseApiPath := os.Getenv("IPINFO_BASE_API_PATH")

	return &IpInfoAppConfig{
		Addr:     addr,
		Protocol: protocol,

		MaxHeaderBytes:    maxHeaderBytes,
		ReadTimeout:       time.Duration(readTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(readHeaderTimeout) * time.Second,
		WriteTimeout:      time.Duration(writeTimeout) * time.Second,
		IdleTimeout:       time.Duration(idleTimeout) * time.Second,

		CertFile: certFile,
		KeyFile:  keyFile,

		HandlerConfig: handler.AppHandlerConfig{
			BaseApiPath: baseApiPath,
		},
	}
}

func checkConfig(cfg *IpInfoAppConfig) {

}

func Bootstrap() *IpInfoAppConfig {
	setBootstrapLogger()
	cfg := loadConfig()
	checkConfig(cfg)
	slog.Info("bootstrap", "status", "ok")
	return cfg
}
