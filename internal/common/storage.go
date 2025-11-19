package common

import (
	"context"
	"database/sql"
)

type Storage interface {
	StartUp() error
	ShutDown() error

	UpdateOption(name, value string, ctx context.Context) error
	GetOption(name string, ctx context.Context) (string, error)
	UpdateRirData(rirTableName string, ip_ranges []IpRange, ctx context.Context) error
}

type IpRange struct {
	RirId           int
	CountryCode     string
	IpVersionId     int
	StartIp         string
	EndIp           string
	Quantity        uint64
	StatusId        int
	StatusChangedAt sql.NullTime
}
