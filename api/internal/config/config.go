// Package config は環境変数からアプリケーション設定を読み込む。
package config

import (
	"errors"
	"os"
)

// ErrMissingDSN は DB_DSN が未設定のときに返る。
var ErrMissingDSN = errors.New("config: DB_DSN is not set")

// DB は DB 接続に必要な設定を保持する。
type DB struct {
	DSN string
}

// LoadDB は環境変数 DB_DSN から DB 設定を読み込む。
// 未設定なら ErrMissingDSN を返す（フォールバックしない）。
func LoadDB() (DB, error) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		return DB{}, ErrMissingDSN
	}
	return DB{DSN: dsn}, nil
}
