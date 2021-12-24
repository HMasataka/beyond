package repository

import (
	"context"

	"github.com/caravan-inc/fankey-server/domain/model"
	"github.com/caravan-inc/fankey-server/domain/repository"
	"github.com/caravan-inc/fankey-server/transactor"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type UserRepository struct {
	connectionProvider transactor.ConnectionProvider
}

func NewUserRepository(connectionProvider transactor.ConnectionProvider) repository.UserRepository {
	return &UserRepository{
		connectionProvider: connectionProvider,
	}
}

func (r *UserRepository) Insert(ctx context.Context, target *model.User) error {
	client := r.connectionProvider.CurrentConnection(ctx)
	return target.Insert(ctx, client, boil.Infer())
}

func (r *UserRepository) Find(ctx context.Context, userID string) (*model.User, error) {
	client := r.connectionProvider.CurrentConnection(ctx)
	return model.FindUser(ctx, client, userID)
}
