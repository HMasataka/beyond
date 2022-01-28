package repository

import (
	"context"

	"github.com/HMasataka/beyond/domain/model"
	"github.com/HMasataka/beyond/domain/repository"
	"github.com/HMasataka/beyond/transactor"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type CardRepository struct {
	connectionProvider transactor.ConnectionProvider
}

func NewCardRepository(connectionProvider transactor.ConnectionProvider) repository.CardRepository {
	return &CardRepository{
		connectionProvider: connectionProvider,
	}
}

func (r *CardRepository) Insert(ctx context.Context, target *model.Card) error {
	client := r.connectionProvider.CurrentConnection(ctx)
	return target.Insert(ctx, client, boil.Infer())
}

func (r *CardRepository) Find(ctx context.Context, userID string) (*model.Card, error) {
	client := r.connectionProvider.CurrentConnection(ctx)
	return model.FindCard(ctx, client, userID)
}
