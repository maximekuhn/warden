package logger

import (
	"context"
	"log/slog"

	"github.com/maximekuhn/warden/internal/auth"
)

// UpgradeLoggerWithRequestId accepts the request context and a parent logger.
// It will return a new logger with the field requestId filled.
// If the request Id is not found in the request context, it will be set to the
// string 'unknown'.
//
// NOTE: the key is required to avoid import cycles in middlewares package.
// When calling this function, set it to middlewares.RequestIdKey
func UpgradeLoggerWithRequestId(reqCtx context.Context, key interface{}, parentLogger *slog.Logger) *slog.Logger {
	reqId, ok := reqCtx.Value(key).(string)
	if !ok {
		reqId = "unknown"
	}
	return parentLogger.With(slog.String("requestId", reqId))
}

// UpgradeLoggerWithRequestId accepts the request context and a parent logger.
// It will return a new logger with the field userId filled.
// If the logged user is not found in the request context, it will be set to the
// string 'unknown'.
//
// NOTE: the key is required to avoid import cycles in middlewares package.
// When calling this function, set it to middlewares.LoggedUserKey
func UpgradeLoggerWithUserId(reqCtx context.Context, key interface{}, parentLogger *slog.Logger) *slog.Logger {
	loggedUser, ok := reqCtx.Value(key).(auth.User)
	userID := "unknown"
	if ok {
		userID = loggedUser.ID.String()
	}
	return parentLogger.With(slog.String("userId", userID))

}
