// Package di はアプリケーションの合成ルートを提供する。
// 各層の依存をここだけで結線する（ADR 0004）。
package di

import (
	"context"
	"fmt"

	"github.com/HMasataka/beyond/internal/config"
	"github.com/HMasataka/beyond/internal/handler"
	"github.com/HMasataka/beyond/internal/infra/database"
)

// New は依存を結線し、集約 Server と解放用の cleanup を返す。
// cleanup は呼び出し側が停止後に呼ぶ（DB 接続を閉じる）。
func New(ctx context.Context) (*handler.Server, func() error, error) {
	cfg, err := config.LoadDB()
	if err != nil {
		return nil, nil, fmt.Errorf("load db config: %w", err)
	}

	db, err := database.Open(ctx, cfg.DSN)
	if err != nil {
		return nil, nil, fmt.Errorf("open database: %w", err)
	}

	server := handler.NewServer(db)
	cleanup := func() error {
		return db.Close()
	}

	return server, cleanup, nil
}
