package app

import (
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

var Rirs = []string{
	"apnic",
	// "arin",
	// "iana",
	// "lacnic",
	// "ripencc",
}

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
	if err := p.database.StartUp(); err != nil {
		return err
	}
	if err := p.cache.StartUp(); err != nil {
		p.database.ShutDown()
		return err
	}

	go p.ShutDownHandler()

	return func() error {
		for {
			if err := p.Update(); err != nil {
				return err
			}
			time.Sleep(100 * time.Minute)
			slog.Info("successfuly updated")
		}
	}()
}

func (p *IpInfoUpdaterApp) Update() error {
	var wg sync.WaitGroup

	// var tableName string
	// existsTableName, err := p.database.GetCurrentIpRangesName()
	// if err != nil {
	// 	return err
	// }
	// if existsTableName == common.AIpRangesTable {
	// 	tableName = common.BIpRangesTable
	// } else {
	// 	tableName = common.AIpRangesTable
	// }

	wg.Add(len(Rirs))
	for _, rirName := range Rirs {
		go func() {
			defer func() {
				slog.Info("finish", "rir", rirName)
				wg.Done()
			}()
			rir := NewRir(rirName, p.config.RegistryFilePath)
			err := rir.Update(p.database, "ip_ranges_a")
			if err != nil {
				slog.Error("dowload rir", "rir", rirName, "error", err)
				return
			}

		}()
	}

	wg.Wait()

	return nil
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
