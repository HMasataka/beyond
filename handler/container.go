package handler

import (
	"sync"
)

type HandlerContainer struct {
	HealthHandler
}

var (
	container *HandlerContainer
	once      sync.Once
)

func NewHandlerOnce() *HandlerContainer {
	once.Do(func() {
		container = newContainer()
	})

	return container
}

func newContainer() *HandlerContainer {
	return &HandlerContainer{
		HealthHandler: NewHealthHandler(),
	}
}
