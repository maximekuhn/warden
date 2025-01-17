package middlewares

import (
	"log/slog"
	"net/http"
	"time"
)

type LoggerMiddleware struct {
	logger *slog.Logger
}

func NewLoggerMiddleware(l *slog.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{logger: l}
}

func (mw *LoggerMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId, ok := r.Context().Value(RequestIdKey).(string)
		if !ok {
			reqId = "unknown"
		}
		l := mw.logger.With(
			slog.String("requestId", reqId),
			slog.String("method", r.Method),
			slog.String("uri", r.RequestURI))

		start := time.Now()
		l.Info("start", slog.Int64("timestamp", start.Unix()))

		// wrapper to get back response status code
		rw := responseWriterWithStatusCode{
			ResponseWriter: w,
			code:           http.StatusOK, // default status code
		}
		next.ServeHTTP(&rw, r)

		end := time.Now()
		duration := end.Sub(start).String()
		code := rw.code

		l.Info("end",
			slog.Int64("timestamp", end.Unix()),
			slog.String("duration", duration),
			slog.Int("code", code))
	})
}

type responseWriterWithStatusCode struct {
	http.ResponseWriter
	code int
}

func (r *responseWriterWithStatusCode) WriteHeader(statusCode int) {
	r.code = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
