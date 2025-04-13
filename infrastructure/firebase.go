package infrastructure

import (
	"context"
	"sync"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/messaging"
	"github.com/HMasataka/beyond/config"
	"google.golang.org/api/option"
)

var (
	firebaseAuthClient      *auth.Client
	firebaseMessagingClient *messaging.Client
	firebaseAuthOnce        sync.Once
	firebaseMessagingOnce   sync.Once
)

type TokenVerifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
	VerifyIDTokenAndCheckRevoked(ctx context.Context, idToken string) (*auth.Token, error)
}

func NewFirebaseAuthClient(ctx context.Context) TokenVerifier {
	firebaseAuthOnce.Do(func() {
		firebaseApp, err := firebase.NewApp(ctx, nil)
		if err != nil {
			panic(err)
		}

		client, err := firebaseApp.Auth(ctx)
		if err != nil {
			panic(err)
		}

		firebaseAuthClient = client
	})

	return firebaseAuthClient
}

func NewFirebaseMessagingClient(ctx context.Context, cfg *config.Config) *messaging.Client {
	firebaseMessagingOnce.Do(func() {
		firebaseApp, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(cfg.Firebase.Credentials))
		if err != nil {
			panic(err)
		}

		client, err := firebaseApp.Messaging(ctx)
		if err != nil {
			panic(err)
		}

		firebaseMessagingClient = client
	})

	return firebaseMessagingClient
}
