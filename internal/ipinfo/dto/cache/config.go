package cache

import (
	"errors"
	"fmt"
	"os"

	"github.com/KeilWin/ipinfo/internal/common"
)

const componentName = "CACHE"

type CacheConfig struct {
	common.Config

	BasePrefix string

	Type CacheType
}

func (p *CacheConfig) NewVariableName(name string) string {
	return fmt.Sprintf("%s_%s", p.BasePrefix, name)
}

func (p *CacheConfig) Load() error {
	var hasError bool

	cacheTypeName := p.NewVariableName("TYPE")
	p.Type = CacheType(os.Getenv(cacheTypeName))

	if hasError {
		return errors.New("loading cache config")
	}
	return nil
}

func (p *CacheConfig) Check() error {
	return nil
}

func NewCacheConfig(appPrefix string) *CacheConfig {
	return &CacheConfig{
		BasePrefix: common.NewBasePrefix(appPrefix, componentName),
	}
}
