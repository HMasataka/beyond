package context

import "context"

const ctxUserIDKey = "UserID"

func GetUserID(ctx context.Context) string {
	return ctx.Value(ctxUserIDKey).(string)
}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ctxUserIDKey, userID)
}

func ExistUserID(ctx context.Context) bool {
	_, ok := ctx.Value(ctxUserIDKey).(string)
	return ok
}
