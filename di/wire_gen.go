// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"context"
	"github.com/HMasataka/beyond/config"
	"github.com/HMasataka/beyond/handler"
	"github.com/HMasataka/beyond/infrastructure"
	"github.com/HMasataka/beyond/infrastructure/driver"
	"github.com/HMasataka/beyond/middleware"
	"github.com/HMasataka/beyond/usecase"
	"github.com/HMasataka/transactor/rdbms"
)

// Injectors from auth-middleware.wire.go:

func InitializeAuthMiddlewareHandler(ctx context.Context, cfg *config.Config) *middleware.Auth {
	conn := infrastructure.NewConnectionOnce(cfg)
	connectionProvider := rdbms.NewConnectionProvider(conn)
	transactor := rdbms.NewTransactor(connectionProvider)
	tokenVerifier := infrastructure.NewFirebaseAuthClient(ctx)
	clientProvider := rdbms.NewClientProvider(connectionProvider)
	driverContainer := driver.NewDriverOnce(clientProvider)
	auth := middleware.NewAuthMiddleware(transactor, tokenVerifier, driverContainer)
	return auth
}

// Injectors from server.wire.go:

func InitializeServerHandler(cfg *config.Config) *handler.HandlerContainer {
	conn := infrastructure.NewConnectionOnce(cfg)
	connectionProvider := rdbms.NewConnectionProvider(conn)
	transactor := rdbms.NewTransactor(connectionProvider)
	clientProvider := rdbms.NewClientProvider(connectionProvider)
	driverContainer := driver.NewDriverOnce(clientProvider)
	useCaseContainer := usecase.NewUseCaseOnce(transactor, driverContainer)
	handlerContainer := handler.NewHandlerOnce(useCaseContainer)
	return handlerContainer
}
