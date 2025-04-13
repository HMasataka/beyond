package middleware

import (
	"net/http"
	"time"

	ctxutil "github.com/caravan-inc/oshi-card-card-recommender/context"
)

func WithCurrentTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := ctxutil.WithRequestTime(r.Context(), time.Now())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
