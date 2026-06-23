// Package database は MySQL への *sql.DB 接続を構築する。
package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// database/sql に mysql ドライバを登録するための副作用 import。
	_ "github.com/go-sql-driver/mysql"
)

const (
	maxOpenConns    = 25
	maxIdleConns    = 25
	connMaxLifetime = 5 * time.Minute
	connMaxIdleTime = 5 * time.Minute
	pingTimeout     = 5 * time.Second
)

// Open は DSN から MySQL 接続を開き、プールを設定して起動時 ping を行う。
// ping に失敗したら接続を閉じてエラーを返す（fail-fast）。
func Open(ctx context.Context, dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetConnMaxIdleTime(connMaxIdleTime)

	pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping mysql: %w", err)
	}

	return db, nil
}
