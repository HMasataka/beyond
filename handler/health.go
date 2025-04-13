package handler

import (
	"net/http"

	"github.com/caravan-inc/oshi-card-card-recommender/handler/payload"
)

type HealthHandler interface {
	Healthz(w http.ResponseWriter, r *http.Request) error
}

type Health struct{}

func NewHealthHandler() HealthHandler {
	return &Health{}
}

// @Summary Healthz
// @Description ヘルスチェック
// @ID Healthz
// @Tags health
// @Produce json
// @Success 200 {object} payload.GetHealthResponse
// @Failure 400 {object} payload.Error
// @Router /healthz [get]
func (u *Health) Healthz(w http.ResponseWriter, r *http.Request) error {
	return writeResponse(w, &payload.GetHealthResponse{})
}
