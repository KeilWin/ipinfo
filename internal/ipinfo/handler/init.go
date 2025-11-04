package handler

import (
	"fmt"
	"log/slog"
	"net/http"
)

func initHandler(handler *http.ServeMux, handlerConfig *HandlerConfig) {
	healthPath := fmt.Sprintf("GET %s/health", handlerConfig.ApiBasePath)
	handler.Handle(healthPath, newHealthHandler())
	slog.Info("added health path", "path", healthPath)

	ipv4Path := fmt.Sprintf("GET %s/ipv4/{ipAddress}", handlerConfig.ApiBasePath)
	handler.Handle(ipv4Path, NewIpV4Handler())
	slog.Info("added ipv4 path", "path", ipv4Path)

	ipv6Path := fmt.Sprintf("GET %s/ipv6/{ipAddress}", handlerConfig.ApiBasePath)
	handler.Handle(ipv6Path, NewIpV6Handler())
	slog.Info("added ipv6 path", "path", ipv6Path)
}

func NewAppHandler(handlerConfig *HandlerConfig) *http.ServeMux {
	handler := http.NewServeMux()
	initHandler(handler, handlerConfig)
	return handler
}
