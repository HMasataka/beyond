package repository

import (
	"context"

	"github.com/caravan-inc/fankey-server/domain/model"
)

type UserRepository interface {
	Insert(ctx context.Context, target *model.User) error
	Find(ctx context.Context, userID string) (*model.User, error)
}
