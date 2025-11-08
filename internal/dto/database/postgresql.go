package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/KeilWin/ipinfo/internal/common"
	"github.com/lib/pq"
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
	slog.Info("database connected")

	p.Db.SetMaxOpenConns(p.Config.MaxOpenConnections)
	p.Db.SetMaxIdleConns(p.Config.MaxIdleConnections)

	p.Db.SetConnMaxLifetime(p.Config.ConnectionMaxLifetime * 60 * time.Second)
	p.Db.SetConnMaxIdleTime(p.Config.ConnectionMaxIdleTime * 60 * time.Second)

	if err = p.Db.Ping(); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}
	return err
}

func (p *PostgreSqlDatabase) ShutDown() error {
	return p.Db.Close()
}

func (p *PostgreSqlDatabase) GetIpInfo(ipAddress string) string {
	return "123.123.123.123"
}

func (p *PostgreSqlDatabase) GetCurrentIpRangesName() (string, error) {
	var aExists, bExists bool
	err := p.Db.QueryRow(`SELECT EXISTS (
            SELECT FROM information_schema.tables 
            WHERE table_schema = 'public' 
            AND table_name = '$1'
        );`, common.AIpRangesTable).Scan(&aExists)
	if err != nil {
		return "", nil
	}
	err = p.Db.QueryRow(`SELECT EXISTS (
            SELECT FROM information_schema.tables 
            WHERE table_schema = 'public' 
            AND table_name = '$1'
        )`, common.BIpRangesTable).Scan(&bExists)
	if err != nil {
		return "", nil
	}

	if aExists && bExists {
		return "", errors.New("two ip_range tables exists")
	} else if !aExists && !bExists {
		return "", errors.New("no one ip_range tables exists")
	}

	if aExists {
		return common.AIpRangesTable, nil
	}
	return common.BIpRangesTable, nil
}

func (p *PostgreSqlDatabase) CopyToIpRangesFromArray(table string, ip_ranges []common.IpRange) error {
	slog.Info("start copy to database", "table", table, "quantity", len(ip_ranges))
	tx, err := p.Db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(pq.CopyIn(table, "rir", "country_code", "version_ip", "start_ip", "end_ip", "status", "created_at"))
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}

	for i, ip_range := range ip_ranges {
		_, err = stmt.Exec(ip_range.Rir, ip_range.CountryCode, ip_range.VersionIp, ip_range.StartIp, ip_range.EndIp, ip_range.Status, ip_range.CreatedAt)
		if err != nil {
			return fmt.Errorf("exec[%d] = '%v': %w", i, ip_range, err)
		}
	}

	if _, err = stmt.Exec(); err != nil {
		return fmt.Errorf("finish copy: %w", err)
	}

	if err = stmt.Close(); err != nil {
		return fmt.Errorf("stmt close: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

func NewPostgreSqlDatabase(cfg *DatabaseConfig) *PostgreSqlDatabase {
	return &PostgreSqlDatabase{
		Config: cfg,
	}
}
