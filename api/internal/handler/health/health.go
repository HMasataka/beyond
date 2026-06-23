package health

import (
	"context"

	"github.com/HMasataka/beyond/internal/openapi"
)

// Health は health 機能のオペレーションを実装する。
type Health struct{}

// New は health ハンドラを生成する。
func New() *Health {
	return &Health{}
}

// GetHealthz はサービスの稼働を示す 200 応答を返す。
func (h *Health) GetHealthz(_ context.Context, _ openapi.GetHealthzRequestObject) (openapi.GetHealthzResponseObject, error) {
	return openapi.GetHealthz200JSONResponse{Status: openapi.HealthStatusStatusOk}, nil
}
