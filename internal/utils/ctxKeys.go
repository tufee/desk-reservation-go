package utils

import "golang.org/x/net/context"

type ctxKey string

const (
	CreateUserKey        ctxKey = "CreateUser"
	CreateReservationKey ctxKey = "CreateReservation"
)

func SetContextValue[T any](ctx context.Context, key ctxKey, value T) context.Context {
	return context.WithValue(ctx, key, value)
}

func GetContextValue[T any](ctx context.Context, key ctxKey) (T, bool) {
	val := ctx.Value(key)
	if val == nil {
		var zero T
		return zero, false
	}

	casted, ok := val.(T)
	return casted, ok
}
