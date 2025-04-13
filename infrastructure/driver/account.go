package driver

import (
	"context"

	"github.com/HMasataka/beyond/domain/entity"
)

func (r *accountDriver) FindByFirebaseUID(ctx context.Context, firebaseUID string) (*entity.Account, error) {
	client := r.clientProvider.CurrentClient(ctx)

	account, err := entity.Accounts(entity.AccountWhere.FirebaseID.EQ(firebaseUID)).One(ctx, client)
	if wrapError(err) != nil {
		return nil, err
	}

	return account, nil
}
