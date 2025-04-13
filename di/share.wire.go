//go:build wireinject
// +build wireinject

package di

import (
	"github.com/HMasataka/beyond/infrastructure"
	"github.com/HMasataka/transactor/rdbms"
	"github.com/google/wire"
)

var DatabaseClient = wire.NewSet(
	infrastructure.NewConnectionOnce,
	rdbms.NewConnectionProvider,
	rdbms.NewClientProvider,
	rdbms.NewTransactor,
)
