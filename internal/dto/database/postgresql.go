package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/KeilWin/ipinfo/internal/common"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type IpAddressInfoRow struct {
	Id               string `json:"id"`
	RirName          string `json:"rirName"`
	CountryCode      string `json:"countryCode"`
	IpAddressVersion string `json:"ipAddressVersion"`
	IpRangeStart     string `json:"ipRangeStart"`
	IpRangeEnd       string `json:"ipRangeEnd"`
	IpRangeQuantity  string `json:"ipRangeQuantity"`
	Status           string `json:"status"`
	StatusUpdatedAt  string `json:"statusUpdatedAt"`
}

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

func (p *PostgreSqlDatabase) GetIpInfo(ipAddress string) (*IpAddressInfoRow, error) {
	var ipInfoRow IpAddressInfoRow
	err := p.Db.QueryRow("SELECT * FROM ip_ranges WHERE start_ip <= '$1'::inet AND end_ip > '$2'::inet LIMIT 1", ipAddress, ipAddress).Scan(&ipInfoRow)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &ipInfoRow, nil
}

func (p *PostgreSqlDatabase) UpdateOption(name, value string, ctx context.Context) error {
	_, err := p.Db.ExecContext(ctx, `INSERT INTO options (name, value) VALUES ($1, $2) 
	ON CONFLICT (name) DO 
	UPDATE SET value = EXCLUDED.value;`, name, value)
	if err != nil {
		return fmt.Errorf("insert with on conflict: %w", err)
	}
	return nil
}

func (p *PostgreSqlDatabase) GetOption(name string, ctx context.Context) (string, error) {
	var result string
	err := p.Db.QueryRowContext(ctx, "SELECT value FROM options WHERE name=$1 LIMIT 1", name).Scan(&result)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (p *PostgreSqlDatabase) UpdateRirData(rirTableName string, ip_ranges []common.IpRange, ctx context.Context) error {
	tx, err := p.Db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY", rirTableName))
	if err != nil {
		return fmt.Errorf("truncate: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, pq.CopyIn(rirTableName, "country_code", "ip_version_id", "start_ip", "end_ip", "quantity", "status_id", "status_changed_at"))
	if err != nil {
		return fmt.Errorf("stmt open: %w", err)
	}

	for i, ip_range := range ip_ranges {
		_, err = stmt.ExecContext(ctx, ip_range.CountryCode, ip_range.IpVersionId, ip_range.StartIp, ip_range.EndIp, ip_range.Quantity, ip_range.StatusId, ip_range.StatusChangedAt)
		if err != nil {
			return fmt.Errorf("exec[%d] = '%v': %w", i, ip_range, err)
		}
	}

	if _, err = stmt.ExecContext(ctx); err != nil {
		return fmt.Errorf("finish copy: %w", err)
	}

	if err = stmt.Close(); err != nil {
		return fmt.Errorf("stmt close: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	_, err = p.Db.ExecContext(ctx, "REFRESH MATERIALIZED VIEW ip_ranges;")
	if err != nil {
		return err
	}
	return nil
}

func NewPostgreSqlDatabase(cfg *DatabaseConfig) *PostgreSqlDatabase {
	return &PostgreSqlDatabase{
		Config: cfg,
	}
}
