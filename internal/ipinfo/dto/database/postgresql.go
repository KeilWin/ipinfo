package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type PostgreSqlDatabase struct {
	Database

	Config *DatabaseConfig
	Db     *sql.DB
}

func (p *PostgreSqlDatabase) StartUp() error {
	var err error

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		p.Config.Host, p.Config.Port, p.Config.User, p.Config.Password, p.Config.Name,
	)

	p.Db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	p.Db.SetMaxOpenConns(p.Config.MaxOpenConnections)
	p.Db.SetMaxIdleConns(p.Config.MaxIdleConnections)

	p.Db.SetConnMaxLifetime(p.Config.ConnectionMaxLifetime * 60 * time.Second)
	p.Db.SetConnMaxIdleTime(p.Config.ConnectionMaxIdleTime * 60 * time.Second)

	return p.Db.Ping()
}

func (p *PostgreSqlDatabase) ShutDown() error {
	return p.Db.Close()
}

func (p *PostgreSqlDatabase) GetIpInfo(ipAddress string) string {
	return "123.123.123.123"
}

func NewPostgreSqlDatabase(cfg *DatabaseConfig) *PostgreSqlDatabase {
	return &PostgreSqlDatabase{
		Config: cfg,
	}
}
