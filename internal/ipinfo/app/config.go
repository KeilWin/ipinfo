package app

import (
	"errors"
	"fmt"
	"strings"

	"github.com/KeilWin/ipinfo/internal/common"
	"github.com/KeilWin/ipinfo/internal/ipinfo/dto/cache"
	"github.com/KeilWin/ipinfo/internal/ipinfo/dto/database"
	"github.com/KeilWin/ipinfo/internal/ipinfo/handler"
)

const AppName = "IPINFO"

func dotEnvFilename() string {
	return fmt.Sprintf("./configs/%s.env", strings.ToLower(AppName))
}

type IpInfoAppConfig struct {
	common.Config
	BasePrefix string
	Server     *ServerConfig
	Handler    *handler.HandlerConfig
	Cache      *cache.CacheConfig
	Database   *database.DatabaseConfig
}

func (p *IpInfoAppConfig) Load() error {
	if p.Server.Load() != nil || p.Handler.Load() != nil || p.Cache.Load() != nil || p.Database.Load() != nil {
		return errors.New("loading app config")
	}
	return nil
}

func (p *IpInfoAppConfig) Check() error {
	if p.Server.Check() != nil || p.Handler.Check() != nil || p.Cache.Check() != nil || p.Database.Check() != nil {
		return errors.New("checking app config")
	}
	return nil
}

func (p *IpInfoAppConfig) Protocol() ProtocolType {
	return p.Server.Protocol
}

func (p *IpInfoAppConfig) CertFile() string {
	return p.Server.CertFile
}

func (p *IpInfoAppConfig) KeyFile() string {
	return p.Server.KeyFile
}

func NewIpInfoAppConfig() *IpInfoAppConfig {
	return &IpInfoAppConfig{
		BasePrefix: AppName,
		Server:     NewServerConfig(AppName),
		Handler:    handler.NewHandlerConfig(AppName),
		Cache:      cache.NewCacheConfig(AppName),
		Database:   database.NewDatabaseConfig(AppName),
	}
}
