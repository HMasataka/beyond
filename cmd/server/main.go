package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/caravan-inc/fankey-server/application/api"
	"github.com/caravan-inc/fankey-server/application/api/pet"
	"github.com/caravan-inc/fankey-server/application/api/user"
	middleware2 "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ctxUserID string

const ctxUserIDKey ctxUserID = "user"

func withUserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ctxUserIDKey, "123")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	var port = flag.Int("port", 8080, "Port for test HTTP server")

	flag.Parse()

	petStore := pet.NewPet()
	userStore := user.NewUser()

	swagger, err := api.GetSwagger()
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(withUserID)

	r.Use(middleware2.OapiRequestValidator(swagger))

	r.Group(func(r chi.Router) {
		pet.HandlerFromMux(petStore, r)
		user.HandlerFromMux(userStore, r)
	})

	s := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("0.0.0.0:%d", *port),
	}

	log.Fatal(s.ListenAndServe())
}
