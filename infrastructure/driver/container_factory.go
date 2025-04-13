package driver

import (
	"sync"

	dd "github.com/HMasataka/beyond/domain/driver"
	"github.com/HMasataka/transactor/rdbms"
)

var (
	driverContainer *dd.DriverContainer
	driverOnce      sync.Once
)

func NewDriverOnce(clientProvider rdbms.ClientProvider) *dd.DriverContainer {
	driverOnce.Do(func() {
		driverContainer = newDriver(clientProvider)
	})

	return driverContainer
}

func newDriver(clientProvider rdbms.ClientProvider) *dd.DriverContainer {
	return &dd.DriverContainer{
		AccountDriver: NewAccountDriver(clientProvider),
		UserDriver:    NewUserDriver(clientProvider),
	}
}
