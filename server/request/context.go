package request

import (
	"context"
	"github.com/TheYeung1/yata-server/model"
)

type ctxKey string

const userIDContextKey ctxKey = "UserID"

// WithUserID stores the userID on the returned context.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

// UserID returns the userID stored on the context.
// If the userID was found it will return true, false otherwise.
func UserID(ctx context.Context) (model.UserID, bool) {
	val := ctx.Value(userIDContextKey)
	str, ok := val.(string)
	return model.UserID(str), ok
}
