package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/HMasataka/beyond/config"
	"github.com/HMasataka/beyond/di"
	"github.com/HMasataka/beyond/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	"golang.org/x/exp/slog"
)

// @title Oshi Card Recommendation API
// @version 1.0
// @description oshi card recommendation api
// @host localhost:8080
func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to create new config %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	httpServer := newHTTPServer(ctx, cfg)

	slog.Info(fmt.Sprintf("REST server listening at %v", httpServer.Addr))

	if err := httpServer.ListenAndServe(); err != nil {
		slog.Info(fmt.Sprintf("HTTP server closed: %v", err))
	}
}

func newHTTPServer(ctx context.Context, cfg *config.Config) *http.Server {
	r := newHandler(ctx, cfg)

	return &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.Port),
		Handler: r,
	}
}

func newHandler(ctx context.Context, cfg *config.Config) http.Handler {
	logger := httplog.NewLogger("game-api", httplog.Options{
		JSON: true,
	})

	skipPaths := []string{"/healthz"}

	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(logger, skipPaths))
	r.Use(middleware.WithCurrentTime)

	container := di.InitializeServerHandler(cfg)
	auth := di.InitializeAuthMiddlewareHandler(ctx, cfg)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Content-Range"},
		AllowCredentials: true,
		ExposedHeaders:   []string{"Content-Range"},
	}))

	r.Route("/healthz", func(r chi.Router) {
		r.Get("/", middleware.Wrap(container.IHealthHandler.Healthz))
	})

	r.Route("/account", func(r chi.Router) {
		r.Use(auth.VerifyToken)

		r.Post("/", middleware.Wrap(container.IAccountHandler.PostAccount))
	})

	r.Route("/user", func(r chi.Router) {
		r.Use(auth.CurrentUser)

		r.Get("/{user_id}", middleware.Wrap(container.IUserHandler.GetUser))
		r.Post("/", middleware.Wrap(container.IUserHandler.PostUser))
	})

	return r
}
