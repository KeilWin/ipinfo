package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/KeilWin/ipinfo/internal/dao"
	"github.com/KeilWin/ipinfo/internal/dto/cache"
	"github.com/KeilWin/ipinfo/internal/dto/database"
	"github.com/KeilWin/ipinfo/internal/handler"
	"github.com/KeilWin/ipinfo/internal/logger"
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

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		p.cache.ShutDown()
	}()
	go func() {
		defer wg.Done()
		p.database.ShutDown()
	}()
	wg.Wait()

	os.Exit(int(utils.ExitSuccess))
}

func (p *IpInfoApp) Start() error {
	slog.Info("app starting...")

	if err := p.database.StartUp(); err != nil {
		return err
	}
	if err := p.cache.StartUp(); err != nil {
		p.database.ShutDown()
		return err
	}

	go p.ShutDownHandler()

	switch p.cfg.Protocol() {
	case ProtocolHTTP:
		slog.Info("starting server", "server_protocol", p.cfg.Protocol())
		return p.server.ListenAndServe()
	case ProtocolHTTPS:
		slog.Info("starting server", "server_protocol", p.cfg.Protocol())
		return p.server.ListenAndServeTLS(p.cfg.CertFile(), p.cfg.KeyFile())
	default:
		return fmt.Errorf("protocol not supported: %s", p.cfg.Protocol())
	}
}

func NewApp(appCfg *IpInfoAppConfig) *IpInfoApp {
	logger := logger.NewAppLogger(appCfg.Logger)
	database, err := database.NewDatabase(appCfg.Database)
	utils.CheckAppFatalError(err)
	cache, err := cache.NewCache(appCfg.Cache)
	utils.CheckAppFatalError(err)
	repository := dao.NewIpAddress(database)
	handler := handler.NewAppHandler(appCfg.Handler, repository)
	server := NewAppServer(handler, appCfg.Server)
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
	cfg, err := bootstrap()
	utils.CheckAppFatalError(err)
	app := NewApp(cfg)
	utils.CheckAppFatalError(app.Start())
}
