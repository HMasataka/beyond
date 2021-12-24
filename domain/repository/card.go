package repository

import (
	"context"

	"github.com/caravan-inc/fankey-server/domain/model"
)

type CardRepository interface {
	Insert(ctx context.Context, target *model.Card) error
	Find(ctx context.Context, cardID string) (*model.Card, error)
}
