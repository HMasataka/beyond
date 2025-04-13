package handler

import (
	"net/http"

	appCtx "github.com/HMasataka/beyond/context"
	"github.com/HMasataka/beyond/handler/payload"
	"github.com/HMasataka/beyond/usecase"
)

type IAccountHandler interface {
	PostAccount(w http.ResponseWriter, r *http.Request) error
}

type Account struct {
	accountUsecase usecase.AccountUseCase
}

func NewAccountHandler(
	usecaseContainer *usecase.UseCaseContainer,
) IAccountHandler {
	return &Account{
		accountUsecase: usecaseContainer.AccountUseCase,
	}
}

// @Summary PostAccount
// @Description アカウントを作成する
// @ID PostAccount
// @Tags account
// @Accept json
// @Produce json
// @Success 200 {object} payload.PostAccountResponse
// @Failure 400 {object} payload.Error
// @Router /account [post]
func (u *Account) PostAccount(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	firebaseUID := appCtx.FirebaseUID(ctx)

	userID, err := u.accountUsecase.Insert(ctx, firebaseUID)
	if err != nil {
		return err
	}

	return writeResponse(w, &payload.PostAccountResponse{
		UserID: userID,
	})
}
