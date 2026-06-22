package handler

import (
	"context"

	"github.com/HMasataka/beyond/internal/openapi"
)

// Handler は OpenAPI から生成された StrictServerInterface を実装する。
type Handler struct{}

var _ openapi.StrictServerInterface = (*Handler)(nil)

// New は API ハンドラを生成する。
func New() *Handler {
	return &Handler{}
}

// GetHealthz はサービスの稼働を示す 200 応答を返す。
func (h *Handler) GetHealthz(_ context.Context, _ openapi.GetHealthzRequestObject) (openapi.GetHealthzResponseObject, error) {
	return openapi.GetHealthz200JSONResponse{Status: openapi.Ok}, nil
}
