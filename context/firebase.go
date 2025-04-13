package context

import "context"

const ctxFirebaseUIDKey = "FirebaseUID"

func FirebaseUID(ctx context.Context) string {
	return ctx.Value(ctxFirebaseUIDKey).(string)
}

func WithFirebaseUID(ctx context.Context, uid string) context.Context {
	return context.WithValue(ctx, ctxFirebaseUIDKey, uid)
}

func ExistFirebaseUID(ctx context.Context) bool {
	_, ok := ctx.Value(ctxFirebaseUIDKey).(string)
	return ok
}
