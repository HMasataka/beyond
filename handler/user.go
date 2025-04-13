package handler

import (
	"encoding/json"
	"net/http"

	appCtx "github.com/HMasataka/beyond/context"
	"github.com/HMasataka/beyond/handler/payload"
	"github.com/HMasataka/beyond/usecase"
	"github.com/go-chi/chi/v5"
)

type IUserHandler interface {
	GetUser(w http.ResponseWriter, r *http.Request) error
	PostUser(w http.ResponseWriter, r *http.Request) error
}

type UserHandler struct {
	userUsecase usecase.UserUseCase
}

func NewUserHandler(
	usecaseContainer *usecase.UseCaseContainer,
) IUserHandler {
	return &UserHandler{
		userUsecase: usecaseContainer.UserUseCase,
	}
}

// @Summary GetUser
// @Description ユーザの詳細を取得する
// @ID GetUser
// @Tags user
// @Produce json
// @Param user_id path string true "UserID"
// @Success 200 {object} payload.GetUserResponse
// @Failure 400 {object} payload.Error
// @Router /user/{user_id} [get]
func (u *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := chi.URLParam(r, "user_id")

	user, err := u.userUsecase.Get(ctx, userID)
	if err != nil {
		return err
	}

	return writeResponse(w, &payload.GetUserResponse{
		ID: user.ID,
	})
}

// @Summary PostUser
// @Description ユーザを作成する
// @ID PostUser
// @Tags user
// @Accept json
// @Produce json
// @Param body body payload.PostUserRequest true "JSON request body"
// @Success 200 {object} payload.PostUserResponse
// @Failure 400 {object} payload.Error
// @Router /user [post]
func (u *UserHandler) PostUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	userID := appCtx.GetUserID(ctx)

	var req payload.PostUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	userID, err := u.userUsecase.Insert(ctx, userID, req.Name)
	if err != nil {
		return err
	}

	return writeResponse(w, &payload.PostUserResponse{
		ID: userID,
	})
}
