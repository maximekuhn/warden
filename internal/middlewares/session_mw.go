package middlewares

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/maximekuhn/warden/internal/auth"
)

type LoggedUserContextKey string

const LoggedUserKey = LoggedUserContextKey("loggedUser")

type SessionMiddleware struct {
	logger  *slog.Logger
	service auth.AuthService
}

func NewSessionMiddleware(l *slog.Logger, s auth.AuthService) *SessionMiddleware {
	return &SessionMiddleware{logger: l, service: s}
}

func (mw *SessionMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId, ok := r.Context().Value(RequestIdKey).(string)
		if !ok {
			reqId = "unknown"
		}
		l := mw.logger.With(slog.String("requestId", reqId))

		cookie, err := r.Cookie(auth.CookieName)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				l.Info("no cookie found")
			} else {
				l.Error("failed to get cookie", slog.String("errMsg", err.Error()))
			}

			l.Info("redirecting user to /login")

			// redirect to login page and don't call next handler
			w.Header().Add("HX-Redirect", "/login") // for HTMX callers
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		l.Info("found cookie")
		user, err := mw.service.Authenticate(r.Context(), *cookie)
		if err != nil {
			mw.handleAuthenticateError(w, r, l, err)
			return
		}
		l.Info("user authenticated", slog.String("userId", user.ID.String()))
		ctx := context.WithValue(r.Context(), LoggedUserKey, *user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (mw *SessionMiddleware) handleAuthenticateError(
	_ http.ResponseWriter,
	_ *http.Request,
	l *slog.Logger,
	err error,
) {
	l.Error("failed to authenticate", slog.String("errMsg", err.Error()))
	// TODO
}
