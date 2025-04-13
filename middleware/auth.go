package middleware

import (
	"fmt"
	"net/http"
	"strings"

	appCtx "github.com/HMasataka/beyond/context"
	"github.com/HMasataka/beyond/domain/driver"
	"github.com/HMasataka/beyond/infrastructure"
	"github.com/HMasataka/transactor"
)

func NewAuthMiddleware(tx transactor.Transactor, tokenVerifier infrastructure.TokenVerifier, driverContainer *driver.DriverContainer) *Auth {
	return &Auth{
		tx:            tx,
		tokenVerifier: tokenVerifier,
		accountDriver: driverContainer.AccountDriver,
	}
}

type Auth struct {
	tx            transactor.Transactor
	tokenVerifier infrastructure.TokenVerifier
	accountDriver driver.AccountDriver
}

func (a *Auth) getRequestToken(r *http.Request) (string, error) {
	t := r.Header.Get("Authorization")

	token := strings.Replace(t, "Bearer ", "", 1)
	if len(token) == 0 {
		return "", fmt.Errorf("token is empty")
	}

	return token, nil
}

func (a *Auth) VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		token, err := a.getRequestToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		verified, err := a.tokenVerifier.VerifyIDTokenAndCheckRevoked(ctx, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx = appCtx.WithFirebaseUID(ctx, verified.UID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Auth) CurrentUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var uid string

		if appCtx.ExistFirebaseUID(ctx) {
			uid = appCtx.FirebaseUID(ctx)
		} else {
			token, err := a.getRequestToken(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			verified, err := a.tokenVerifier.VerifyIDTokenAndCheckRevoked(ctx, token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			uid = verified.UID
		}

		account, err := a.accountDriver.FindByFirebaseUID(ctx, uid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		ctx = appCtx.WithUserID(ctx, account.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
