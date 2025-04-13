//go:build wireinject
// +build wireinject

package di

import (
	"github.com/HMasataka/beyond/config"
	"github.com/HMasataka/beyond/handler"
	"github.com/google/wire"
)

func InitializeServerHandler(cfg *config.Config) *handler.HandlerContainer {
	wire.Build(
		handler.NewHandlerOnce,
	)

	return nil
}
