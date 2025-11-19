package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/KeilWin/ipinfo/internal/common"
	"github.com/KeilWin/ipinfo/internal/dto/cache"
	"github.com/KeilWin/ipinfo/internal/dto/database"
	"github.com/KeilWin/ipinfo/internal/logger"
	"github.com/KeilWin/ipinfo/internal/utils"
)

type IpInfoUpdaterApp struct {
	common.App

	config   *IpInfoUpdaterConfig
	logger   *slog.Logger
	database database.Database
	cache    cache.Cache
}

func (p *IpInfoUpdaterApp) ShutDownHandler() {
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

func (p *IpInfoUpdaterApp) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := p.database.StartUp(); err != nil {
		return err
	}
	if err := p.cache.StartUp(); err != nil {
		p.database.ShutDown()
		return err
	}
	go p.ShutDownHandler()

	var wg sync.WaitGroup
	wg.Add(len(Rirs))
	for _, rir := range Rirs {
		go func() {
			defer func() {
				slog.Info("finish", "rir", rir.DbName)
				wg.Done()
			}()
			rirManager := NewRirManager(rir, p.database, ctx, time.Date(0, 0, 0, 4, 0, 0, 0, time.UTC))
			workLoop := NewWorkLoop(rirManager, 30*time.Minute)
			workLoop()
		}()
	}
	wg.Wait()

	return nil
}

func NewWorkLoop(rirManager *RirManager, retryPause time.Duration) func() {
	return func() {
		slog.Info("start workloop", "rir", rirManager.Rir.DbName)
		for {
			err := rirManager.Start()
			if err != nil {
				slog.Error("update rir", "rir", rirManager.Rir.DbName, "error", err, "retry_after(minutes)", retryPause.Minutes())
				time.Sleep(retryPause)
				continue
			}
		}
	}
}

func NewApp(cfg *IpInfoUpdaterConfig) *IpInfoUpdaterApp {
	logger := logger.NewAppLogger(cfg.Logger)
	database, err := database.NewDatabase(cfg.Database)
	utils.CheckAppFatalError(err)
	cache, err := cache.NewCache(cfg.Cache)
	utils.CheckAppFatalError(err)
	return &IpInfoUpdaterApp{
		config:   cfg,
		logger:   logger,
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
