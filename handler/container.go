package handler

import (
	"sync"

	"github.com/HMasataka/beyond/usecase"
)

type HandlerContainer struct {
	IAccountHandler
	IHealthHandler
	IUserHandler
}

var (
	container *HandlerContainer
	once      sync.Once
)

func NewHandlerOnce(usecaseContainer *usecase.UseCaseContainer) *HandlerContainer {
	once.Do(func() {
		container = newContainer(usecaseContainer)
	})

	return container
}

func newContainer(usecaseContainer *usecase.UseCaseContainer) *HandlerContainer {
	return &HandlerContainer{
		IAccountHandler: NewAccountHandler(usecaseContainer),
		IHealthHandler:  NewHealthHandler(),
		IUserHandler:    NewUserHandler(usecaseContainer),
	}
}
