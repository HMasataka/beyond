// Package readiness は依存先への到達性を確認する readiness 機能を実装する。
package readiness

import (
	"context"
	"time"

	"github.com/HMasataka/beyond/internal/openapi"
)

const pingTimeout = 5 * time.Second

// Pinger は依存先への到達性を確認する最小インターフェース。
// *sql.DB がこれを満たし、テストでは fake を注入する。
type Pinger interface {
	PingContext(ctx context.Context) error
}

// Readiness は readiness 機能のオペレーションを実装する。
type Readiness struct {
	pinger Pinger
}

// New は readiness ハンドラを生成する。
func New(pinger Pinger) *Readiness {
	return &Readiness{pinger: pinger}
}

// GetReadyz は依存先 ping に成功すれば 200、失敗すれば 503 を返す。
func (r *Readiness) GetReadyz(ctx context.Context, _ openapi.GetReadyzRequestObject) (openapi.GetReadyzResponseObject, error) {
	pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	if err := r.pinger.PingContext(pingCtx); err != nil {
		return openapi.GetReadyz503JSONResponse{Status: openapi.ReadinessStatusStatusUnavailable}, nil
	}
	return openapi.GetReadyz200JSONResponse{Status: openapi.ReadinessStatusStatusOk}, nil
}
