package handler

import (
	"sync"

	"github.com/HMasataka/beyond/usecase"
)

type HandlerContainer struct {
	HealthHandler
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
		HealthHandler: NewHealthHandler(),
	}
}
