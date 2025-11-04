package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/KeilWin/ipinfo/internal/utils"
	_ "github.com/lib/pq"
)

type PostgreSqlDatabase struct {
	Database
}

func (p *PostgreSqlDatabase) StartUp() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		"localhost", 5432, "postgres", "12345",
	)

	db, err := sql.Open("postgres", psqlInfo)
	utils.CheckAppFatalError(err)
	defer db.Close()

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)

	db.SetConnMaxLifetime(5 * 60 * time.Second)
	db.SetConnMaxIdleTime(5 * 60 * time.Second)

	utils.CheckAppFatalError(db.Ping())
	slog.Info("successfully connected")
}

func (p *PostgreSqlDatabase) ShutDown() {

}

func (p *PostgreSqlDatabase) GetIpInfo(ipAddress string) {

}

func NewPostgreSqlDatabase() *PostgreSqlDatabase {
	return &PostgreSqlDatabase{}
}
