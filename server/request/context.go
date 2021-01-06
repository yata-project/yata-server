package request

import (
	"context"

	"github.com/TheYeung1/yata-server/model"
	"github.com/sirupsen/logrus"
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

const loggerCtxKey ctxKey = "logger"

// Logger returns the logger stored on the context.
// If a logger was not found a new one will be created with an unknown request ID.
func Logger(ctx context.Context) *logrus.Entry {
	e, ok := ctx.Value(loggerCtxKey).(*logrus.Entry)
	if !ok {
		// TODO: We share this next line with middleware/logger.go; consider refactor so we eliminate the risk of de-sync.
		e = logrus.WithField("requestID", "UNKNOWN")
		e.Warn("failed to find logger on the context")
	}
	return e
}

// WithLogger stores the logger on the returned context.
// Should only be used by the RequestLogger middleware and tests.
func WithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}
