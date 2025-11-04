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
	"github.com/KeilWin/ipinfo/internal/utils"
)

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

	os.Exit(int(utils.ExitSuccess))
}

func (p *IpInfoApp) Start() {
	go p.database.StartUp()
	go p.cache.StartUp()

	go p.ShutDownHandler()

	switch p.cfg.Protocol() {
	case ProtocolHTTP:
		slog.Info("starting server", slog.String("server_protocol", string(p.cfg.Protocol())))
		slog.Error("server error", "error", p.server.ListenAndServe())
	case ProtocolHTTPS:
		slog.Info("starting server", slog.String("server_protocol", string(p.cfg.Protocol())))
		slog.Error("server error", "error", p.server.ListenAndServeTLS(p.cfg.CertFile(), p.cfg.KeyFile()))
	default:
		slog.Error("protocol not supported", slog.String("server_protocol", string(p.cfg.Protocol())))
	}
}

func newIpInfoApp(appCfg *IpInfoAppConfig) *IpInfoApp {
	logger := NewAppLogger(appCfg)
	handler := handler.NewAppHandler(appCfg.Handler)
	server := NewAppServer(handler, appCfg.Server)
	database, err := database.NewDatabase(appCfg.Database)
	utils.CheckAppFatalError(err)
	cache, err := cache.NewCache(appCfg.Cache)
	utils.CheckAppFatalError(err)
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
	defer func() {
		if r := recover(); r != nil {
			slog.Error("app fatal error", "panic", r)
		}
	}()
	appCfg, err := bootstrap()
	utils.CheckAppFatalError(err)
	ipInfoApp := newIpInfoApp(appCfg)
	slog.Info("create app", "status", "ok")
	ipInfoApp.Start()
}
