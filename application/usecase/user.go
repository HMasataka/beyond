package usecase

import (
	"context"

	"github.com/caravan-inc/fankey-server/domain/model"
	"github.com/caravan-inc/fankey-server/domain/model/inventory"
	"github.com/caravan-inc/fankey-server/domain/repository"
	"github.com/caravan-inc/fankey-server/transactor"
)

type UserUseCase interface {
	Insert(ctx context.Context, userID, name string) error
	Find(ctx context.Context, userID string) (inventory.User, error)
}

type userUseCase struct {
	transactor     transactor.Transactor
	userRepository repository.UserRepository
}

func NewUserUseCase(transactor transactor.Transactor, userRepository repository.UserRepository) UserUseCase {
	return &userUseCase{
		transactor:     transactor,
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
