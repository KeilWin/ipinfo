package app

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/KeilWin/ipinfo/internal/common"
	"github.com/KeilWin/ipinfo/internal/dto/cache"
	"github.com/KeilWin/ipinfo/internal/dto/database"
	"github.com/KeilWin/ipinfo/internal/logger"
	"github.com/KeilWin/ipinfo/internal/utils"
)

const AppName = "IPINFO_UPDATER"

type DurationType string

const (
	Second = "second"
	Minute = "minute"
	Hour   = "hour"
)

func dotEnvFilename() string {
	return fmt.Sprintf("./configs/%s.env", strings.ToLower(AppName))
}

type IpInfoUpdaterConfig struct {
	common.Config

	BasePrefix string

	Logger   *logger.LoggerConfig
	Database *database.DatabaseConfig
	Cache    *cache.CacheConfig

	RegistryFilePath string
	DurationType     DurationType
	UpdateFrequency  time.Duration
}

func (p *IpInfoUpdaterConfig) Load() error {
	var err error
	var hasError bool

	hasError = p.Logger.Load() != nil || p.Database.Load() != nil || p.Cache.Load() != nil

	registryFilepathName := p.NewVariableName("REGISTRY_FILEPATH")
	p.RegistryFilePath = os.Getenv(registryFilepathName)

	durationTypeName := p.NewVariableName("DURATION_TYPE")
	p.DurationType = DurationType(os.Getenv(durationTypeName))

	duration := func() time.Duration {
		switch p.DurationType {
		case Second:
			return time.Second
		case Minute:
			return time.Minute
		case Hour:
			return time.Hour
		default:
			return time.Hour
		}
	}()

	updateFrequencyName := p.NewVariableName("UPDATE_FREQUENCY")
	updateFrequency, err := strconv.Atoi(os.Getenv(updateFrequencyName))
	hasError = CheckLoadConfigError(err, updateFrequencyName)
	p.UpdateFrequency = time.Duration(updateFrequency) * duration

	if hasError {
		return errors.New("loading app config")
	}
	return nil
}

func (p *IpInfoUpdaterConfig) NewVariableName(name string) string {
	return fmt.Sprintf("%s_%s", p.BasePrefix, name)
}

func (p *IpInfoUpdaterConfig) Check() error {
	return nil
}

func CheckLoadConfigError(err error, name string) bool {
	return utils.CheckLoadConfigError(err, name, AppName)
}

func NewIpInfoUpdaterConfig() *IpInfoUpdaterConfig {
	return &IpInfoUpdaterConfig{
		BasePrefix: AppName,

		Logger:   logger.NewLoggerConfig(),
		Cache:    cache.NewCacheConfig(AppName),
		Database: database.NewDatabaseConfig(AppName),
	}
}
