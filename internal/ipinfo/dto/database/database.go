package database

import (
	"errors"
	"fmt"
)

type DatabaseType string

const (
	PostgreSqlDatabaseType DatabaseType = "postgresql"
	ClickHouseDatabaseType DatabaseType = "clickhouse"
)

type Database interface {
	StartUp()
	ShutDown()

	GetIpInfo(ipAddress string)
}

func NewDatabase(databaseType DatabaseType) (Database, error) {
	switch databaseType {
	case PostgreSqlDatabaseType:
		return NewPostgreSqlDatabase(), nil
	case ClickHouseDatabaseType:
		return nil, errors.New("clickhouse not implemented")
	default:
		return nil, fmt.Errorf("unknown database type: %s", databaseType)
	}
}
