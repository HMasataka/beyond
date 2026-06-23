package handler

import (
	"github.com/HMasataka/beyond/internal/handler/health"
	"github.com/HMasataka/beyond/internal/handler/readiness"
	"github.com/HMasataka/beyond/internal/openapi"
)

// Server は機能別ハンドラを集約し、生成された StrictServerInterface を満たす。
// 機能ハンドラの型名は機能ごとに一意にする（匿名フィールド名が型名になり衝突するため）。
type Server struct {
	*health.Health
	*readiness.Readiness
}

var _ openapi.StrictServerInterface = (*Server)(nil)

// NewServer は機能別ハンドラを結線して集約する。
func NewServer(pinger readiness.Pinger) *Server {
	return &Server{
		Health:    health.New(),
		Readiness: readiness.New(pinger),
	}
}
