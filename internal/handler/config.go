package handler

import (
	"errors"
	"fmt"
	"os"

	"github.com/KeilWin/ipinfo/internal/common"
)

const componentName = "API"

type HandlerConfig struct {
	common.Config

	BasePrefix string

	ApiBasePath string
}

func (p *HandlerConfig) NewVariableName(name string) string {
	return fmt.Sprintf("%s_%s", p.BasePrefix, name)
}

func (p *HandlerConfig) Load() error {
	var hasError bool

	apiBasePathName := p.NewVariableName("BASE_PATH")
	p.ApiBasePath = os.Getenv(apiBasePathName)

	if hasError {
		return errors.New("loading handler config")
	}
	return nil
}

func (p *HandlerConfig) Check() error {
	return nil
}

func NewHandlerConfig(appPrefix string) *HandlerConfig {
	return &HandlerConfig{
		BasePrefix: common.NewBasePrefix(appPrefix, componentName),
	}
}
