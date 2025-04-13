package usecase

import (
	"sync"

	"github.com/HMasataka/beyond/domain/driver"
	"github.com/HMasataka/transactor"
)

type UseCaseContainer struct {
	UserUseCase UserUseCase
}

var (
	container *UseCaseContainer
	once      sync.Once
)

func NewUseCaseOnce(tx transactor.Transactor, driverContainer *driver.DriverContainer) *UseCaseContainer {
	once.Do(func() {
		container = newContainer(tx, driverContainer)
	})

	return container
}

func newContainer(tx transactor.Transactor, driverContainer *driver.DriverContainer) *UseCaseContainer {
	return &UseCaseContainer{
		UserUseCase: NewUserUseCase(tx, driverContainer),
	}
}
