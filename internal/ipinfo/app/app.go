package app

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type IpInfoApp struct {
	cfg     *IpInfoAppConfig
	logger  *slog.Logger
	handler *http.ServeMux
	server  *http.Server
}

func (p *IpInfoApp) ShutDownHandler() {
	shutDownSignal := make(chan os.Signal, 1)
	signal.Notify(shutDownSignal, syscall.SIGINT, syscall.SIGTERM)
	sig := <-shutDownSignal
	slog.Info("received signal to term", "sig", sig)
	os.Exit(0)
}

func (p *IpInfoApp) Start() {
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
	handler := NewAppHandler(appCfg)
	server := NewAppServer(handler, appCfg)
	return &IpInfoApp{
		cfg:     appCfg,
		logger:  logger,
		handler: handler,
		server:  server,
	}
}

func Start() {
	appCfg := Bootstrap()
	ipInfoApp := newIpInfoApp(appCfg)
	slog.Info("create app", "status", "ok")
	ipInfoApp.Start()
}
