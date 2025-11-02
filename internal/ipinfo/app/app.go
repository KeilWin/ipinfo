package app

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/KeilWin/ipinfo/internal/ipinfo/dto/cache"
	"github.com/KeilWin/ipinfo/internal/ipinfo/dto/database"
	"github.com/KeilWin/ipinfo/internal/ipinfo/handler"
)

type ExitCodeType int

const (
	ExitSuccess ExitCodeType = 0
	ExitError   ExitCodeType = -1
)

func CheckIpInfoAppFatalError(err error) {
	if err != nil {
		slog.Error("app fatal error", "error", err)
		os.Exit(int(ExitError))
	}
}

type IpInfoApp struct {
	cfg      *IpInfoAppConfig
	logger   *slog.Logger
	handler  *http.ServeMux
	server   *http.Server
	database database.Database
	cache    cache.Cache
}

func (p *IpInfoApp) ShutDownHandler() {
	shutDownSignal := make(chan os.Signal, 1)
	signal.Notify(shutDownSignal, syscall.SIGINT, syscall.SIGTERM)

	sig := <-shutDownSignal
	slog.Info("received signal to term", "sig", sig)

	go p.cache.ShutDown()
	go p.database.ShutDown()

	os.Exit(int(ExitSuccess))
}

func (p *IpInfoApp) Start() {
	go p.database.StartUp()
	go p.cache.StartUp()

	go p.ShutDownHandler()

	switch p.cfg.Protocol {
	case ProtocolHTTP:
		slog.Info("starting server", slog.String("server_protocol", string(p.cfg.Protocol)))
		slog.Error("server error", "error", p.server.ListenAndServe())
	case ProtocolHTTPS:
		slog.Info("starting server", slog.String("server_protocol", string(p.cfg.Protocol)))
		slog.Error("server error", "error", p.server.ListenAndServeTLS(p.cfg.CertFile, p.cfg.KeyFile))
	default:
		slog.Error("protocol not supported", slog.String("server_protocol", string(p.cfg.Protocol)))
	}
}

func newIpInfoApp(appCfg *IpInfoAppConfig) *IpInfoApp {
	logger := NewAppLogger(appCfg)
	handler := handler.NewAppHandler(&appCfg.HandlerConfig)
	server := NewAppServer(handler, appCfg)
	database, err := database.NewDatabase(appCfg.DatabaseType)
	CheckIpInfoAppFatalError(err)
	cache, err := cache.NewCache(appCfg.CacheType)
	CheckIpInfoAppFatalError(err)
	return &IpInfoApp{
		cfg:      appCfg,
		logger:   logger,
		handler:  handler,
		server:   server,
		database: database,
		cache:    cache,
	}
}

func Start() {
	appCfg := Bootstrap()
	ipInfoApp := newIpInfoApp(appCfg)
	slog.Info("create app", "status", "ok")
	ipInfoApp.Start()
}
