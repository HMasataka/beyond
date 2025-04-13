//go:build wireinject
// +build wireinject

package di

import (
	"github.com/HMasataka/beyond/config"
	"github.com/HMasataka/beyond/handler"
	"github.com/HMasataka/beyond/infrastructure/driver"
	"github.com/HMasataka/beyond/usecase"
	"github.com/google/wire"
)

func InitializeServerHandler(cfg *config.Config) *handler.HandlerContainer {
	wire.Build(
		DatabaseClient,
		driver.NewDriverOnce,
		usecase.NewUseCaseOnce,
		handler.NewHandlerOnce,
	)

	return nil
}
