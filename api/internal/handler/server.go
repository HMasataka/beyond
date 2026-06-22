package handler

import (
	"github.com/HMasataka/beyond/internal/handler/health"
	"github.com/HMasataka/beyond/internal/openapi"
)

// Server は機能別ハンドラを集約し、生成された StrictServerInterface を満たす。
// 機能を追加するときは internal/handler/<機能>/ を実装し、ここに匿名フィールドとして埋め込む。
// 埋め込みフィールド名は型名になるため、機能ハンドラの型名は機能ごとに一意にする。
type Server struct {
	*health.Health
}

var _ openapi.StrictServerInterface = (*Server)(nil)

// NewServer は機能別ハンドラを結線して集約する。
func NewServer() *Server {
	return &Server{
		Health: health.New(),
	}
}
