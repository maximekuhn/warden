package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type RequestIdContextKey string

const RequestIdKey = RequestIdContextKey("requestId")

type RequestIdMiddleware struct{}

func NewRequestIdMiddleware() *RequestIdMiddleware {
	return &RequestIdMiddleware{}
}

func (mw *RequestIdMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId, ok := r.Context().Value(RequestIdKey).(string)
		if !ok || reqId == "" {
			reqId = uuid.NewString()
		}
		ctx := context.WithValue(r.Context(), RequestIdKey, reqId)
		w.Header().Set("X-Request-Id", reqId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
