package usecase

import (
	"context"

	"github.com/HMasataka/beyond/domain/model"
	"github.com/HMasataka/beyond/domain/model/inventory"
	"github.com/HMasataka/beyond/domain/repository"
	"github.com/HMasataka/beyond/transactor"
)

type UserUseCase interface {
	Insert(ctx context.Context, userID, name string) error
	Find(ctx context.Context, userID string) (inventory.User, error)
}

type userUseCase struct {
	transactor     transactor.Transactor
	userRepository repository.UserRepository
}

func NewUserUseCase(tx transactor.Transactor, userRepository repository.UserRepository) UserUseCase {
	return &userUseCase{
		transactor:     tx,
		userRepository: userRepository,
	}
}

func (u *userUseCase) Insert(ctx context.Context, userID, name string) error {
	user := model.User{
		ID:   userID,
		Name: name,
	}
	err := u.transactor.Required(ctx, func(ctx context.Context) error {
		return u.userRepository.Insert(ctx, &user)
	})

	return err
}

func (u *userUseCase) Find(ctx context.Context, userID string) (inventory.User, error) {
	var user *model.User

	err := u.transactor.Required(ctx, func(ctx context.Context) error {
		var err error
		user, err = u.userRepository.Find(ctx, userID)
		return err
	})

	return inventory.NewUserFrom(user), err
}
