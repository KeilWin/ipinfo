package database

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/KeilWin/ipinfo/internal/common"
	"github.com/KeilWin/ipinfo/internal/utils"
)

type DatabaseType string

const (
	PostgreSqlDatabaseType DatabaseType = "postgresql"
	ClickHouseDatabaseType DatabaseType = "clickhouse"
)

const componentName = "DATABASE"

type DatabaseConfig struct {
	common.Config

	BasePrefix string

	Host     string
	Port     int
	Type     DatabaseType
	User     string
	Password string
	Name     string

	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
	ConnectionMaxIdleTime time.Duration
}

func (p *DatabaseConfig) NewVariableName(name string) string {
	return fmt.Sprintf("%s_%s", p.BasePrefix, name)
}

func (p *DatabaseConfig) Load() error {
	var err error
	var hasError bool

	hostName := p.NewVariableName("HOST")
	p.Host = os.Getenv(hostName)
	portName := p.NewVariableName("PORT")
	p.Port, err = strconv.Atoi(os.Getenv(portName))
	hasError = CheckLoadDatabaseConfigError(err, portName)
	typeName := p.NewVariableName("TYPE")
	p.Type = DatabaseType(os.Getenv(typeName))
	userName := p.NewVariableName("USER")
	p.User = os.Getenv(userName)
	passwordName := p.NewVariableName("PASSWORD")
	p.Password = os.Getenv(passwordName)
	nameName := p.NewVariableName("NAME")
	p.Name = os.Getenv(nameName)

	maxOpenConnectionsName := p.NewVariableName("MAX_OPEN_CONNECTIONS")
	p.MaxOpenConnections, err = strconv.Atoi(os.Getenv(maxOpenConnectionsName))
	hasError = CheckLoadDatabaseConfigError(err, maxOpenConnectionsName)

	maxIdleConnectionsName := p.NewVariableName("MAX_IDLE_CONNECTIONS")
	p.MaxIdleConnections, err = strconv.Atoi(os.Getenv(maxIdleConnectionsName))
	hasError = CheckLoadDatabaseConfigError(err, maxIdleConnectionsName)

	connectionMaxLifetimeName := p.NewVariableName("CONNECTION_MAX_LIFETIME")
	connectionMaxLifetime, err := strconv.Atoi(os.Getenv(connectionMaxLifetimeName))
	hasError = CheckLoadDatabaseConfigError(err, connectionMaxLifetimeName)
	p.ConnectionMaxLifetime = time.Duration(connectionMaxLifetime) * time.Second

	connectionMaxIdleTimeName := p.NewVariableName("CONNECTION_MAX_IDLE_TIME")
	connectionMaxIdleTime, err := strconv.Atoi(os.Getenv(connectionMaxIdleTimeName))
	hasError = CheckLoadDatabaseConfigError(err, connectionMaxIdleTimeName)
	p.ConnectionMaxIdleTime = time.Duration(connectionMaxIdleTime) * time.Second

	if hasError {
		return errors.New("loading database config")
	}
	return nil
}

func (p *DatabaseConfig) Check() error {
	return nil
}

func NewDatabaseConfig(appPrefix string) *DatabaseConfig {
	return &DatabaseConfig{
		BasePrefix: common.NewBasePrefix(appPrefix, componentName),
	}
}

func CheckLoadDatabaseConfigError(err error, name string) bool {
	return utils.CheckLoadConfigError(err, name, componentName)
}
