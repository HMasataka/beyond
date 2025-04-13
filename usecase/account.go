package usecase

import (
	"context"

	"github.com/HMasataka/beyond/domain/driver"
	"github.com/HMasataka/beyond/domain/entity"
	"github.com/HMasataka/beyond/domain/factory"
	"github.com/HMasataka/transactor"
)

type AccountUseCase interface {
	Insert(ctx context.Context, firebaseUID string) (string, error)
}

type accountUseCase struct {
	transactor    transactor.Transactor
	accountDriver driver.AccountDriver
}

func NewAccountUseCase(
	tx transactor.Transactor,
	container *driver.DriverContainer,
) AccountUseCase {
	return &accountUseCase{
		transactor:    tx,
		accountDriver: container.AccountDriver,
	}
}

func (u *accountUseCase) Insert(ctx context.Context, firebaseUID string) (string, error) {
	account := &entity.Account{
		FirebaseID: firebaseUID,
		UserID:     factory.NewID(),
	}

	if err := u.transactor.Required(ctx, func(ctx context.Context) error {
		return u.accountDriver.Insert(ctx, account)
	}); err != nil {
		return "", err
	}

	return account.UserID, nil
}
