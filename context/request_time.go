package context

import (
	"context"
	"time"
)

const ctxRequestTimeKey = "RequestTime"

func RequestTime(ctx context.Context) time.Time {
	return ctx.Value(ctxRequestTimeKey).(time.Time)
}

func WithRequestTime(ctx context.Context, requestTime time.Time) context.Context {
	return context.WithValue(ctx, ctxRequestTimeKey, requestTime)
}

func ExistRequestTime(ctx context.Context) bool {
	_, ok := ctx.Value(ctxRequestTimeKey).(time.Time)
	return ok
}
