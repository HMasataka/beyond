//go:build wireinject
// +build wireinject

package di

import (
	"github.com/caravan-inc/oshi-card-card-recommender/config"
	"github.com/caravan-inc/oshi-card-card-recommender/handler"
	"github.com/google/wire"
)

func InitializeServerHandler(cfg *config.Config) *handler.HandlerContainer {
	wire.Build(
		handler.NewHandlerOnce,
	)

	return nil
}
