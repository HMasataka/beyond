package usecase

import (
	"context"

	"github.com/HMasataka/beyond/domain/driver"
	"github.com/HMasataka/beyond/domain/entity"
	"github.com/HMasataka/transactor"
)

type UserUseCase interface {
	Get(ctx context.Context, firebaseUID string) (*entity.User, error)
	Insert(ctx context.Context, firebaseUID, userID string) (string, error)
}

type userUseCase struct {
	transactor transactor.Transactor
	userDriver driver.UserDriver
}

func NewUserUseCase(
	tx transactor.Transactor,
	container *driver.DriverContainer,
) UserUseCase {
	return &userUseCase{
		transactor: tx,
		userDriver: container.UserDriver,
	}
}

func (u *userUseCase) Get(ctx context.Context, userID string) (*entity.User, error) {
	var user *entity.User

	if err := u.transactor.Required(ctx, func(ctx context.Context) error {
		var err error

		user, err = u.userDriver.Find(ctx, userID)

		return err
	}); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUseCase) Insert(ctx context.Context, userID, name string) (string, error) {
	user := &entity.User{
		ID:   userID,
		Name: name,
	}

	if err := u.transactor.Required(ctx, func(ctx context.Context) error {
		return u.userDriver.Insert(ctx, user)
	}); err != nil {
		return "", err
	}

	return user.ID, nil
}
