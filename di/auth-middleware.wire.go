//go:build wireinject
// +build wireinject

package di

import (
	"context"

	"github.com/HMasataka/beyond/config"
	"github.com/HMasataka/beyond/infrastructure"
	"github.com/HMasataka/beyond/infrastructure/driver"
	"github.com/HMasataka/beyond/middleware"
	"github.com/google/wire"
)

func InitializeAuthMiddlewareHandler(ctx context.Context, cfg *config.Config) *middleware.Auth {
	wire.Build(
		DatabaseClient,
		infrastructure.NewFirebaseAuthClient,
		driver.NewDriverOnce,
		middleware.NewAuthMiddleware,
	)

	return &middleware.Auth{}
}
