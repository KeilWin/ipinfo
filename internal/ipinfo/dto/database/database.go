package database

import (
	"errors"
	"fmt"

	"github.com/KeilWin/ipinfo/internal/common"
)

type Database interface {
	common.Storage

	GetIpInfo(ipAddress string)
}

func NewDatabase(databaseConfig *DatabaseConfig) (Database, error) {
	switch databaseConfig.Type {
	case PostgreSqlDatabaseType:
		return NewPostgreSqlDatabase(), nil
	case ClickHouseDatabaseType:
		return nil, errors.New("clickhouse not implemented")
	default:
		return nil, fmt.Errorf("unknown database type: %s", databaseConfig.Type)
	}
}
