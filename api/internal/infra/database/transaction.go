package database

import (
	"context"
	"database/sql"
	"fmt"
)

// DBTX は repository が tx の有無を意識せずに SQL を発行するための最小実行口。
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// DBProvider は ctx に応じた DBTX を供給する。
type DBProvider interface {
	Executor(ctx context.Context) DBTX
}

type txCtxKey struct{}

// TxManager は tx の有無判定を 1 箇所に集約する。
type TxManager struct {
	db *sql.DB
}

// NewTxManager は TxManager を生成する。
func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{db: db}
}

var _ DBProvider = (*TxManager)(nil)

// Do はネスト時に新たな tx を張らない。内側で commit/rollback が重なるのを避けるため。
func (m *TxManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := ctx.Value(txCtxKey{}).(*sql.Tx); ok {
		return fn(ctx)
	}

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	if err := fn(context.WithValue(ctx, txCtxKey{}, tx)); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback tx: %w (original: %w)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

// Executor は ctx 上の DBTX を返す。
func (m *TxManager) Executor(ctx context.Context) DBTX {
	if tx, ok := ctx.Value(txCtxKey{}).(*sql.Tx); ok {
		return tx
	}
	return m.db
}
