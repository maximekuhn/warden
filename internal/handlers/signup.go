package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/maximekuhn/warden/internal/auth"
	"github.com/maximekuhn/warden/internal/logger"
	"github.com/maximekuhn/warden/internal/middlewares"
	"github.com/maximekuhn/warden/internal/permissions"
	"github.com/maximekuhn/warden/internal/transaction"
	"github.com/maximekuhn/warden/internal/ui/components/errors"
	"github.com/maximekuhn/warden/internal/ui/pages"
	"github.com/maximekuhn/warden/internal/valueobjects"
)

type SignupHandler struct {
	logger      *slog.Logger
	authService *auth.AuthService
	permService *permissions.PermissionsService
	uowProvider transaction.UnitOfWorkProvider
}

func NewSignupHandler(
	l *slog.Logger,
	as *auth.AuthService,
	ps *permissions.PermissionsService,
	uowProvider transaction.UnitOfWorkProvider,
) *SignupHandler {
	return &SignupHandler{
		logger:      l,
		authService: as,
		permService: ps,
		uowProvider: uowProvider,
	}
}

func (h *SignupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.get(w, r)
		return
	}
	if r.Method == http.MethodPost {
		h.post(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *SignupHandler) get(w http.ResponseWriter, r *http.Request) {
	if err := pages.Signup().Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *SignupHandler) post(w http.ResponseWriter, r *http.Request) {
	l := logger.UpgradeLoggerWithRequestId(r.Context(), middlewares.RequestIdKey, h.logger)

	// TODO: handle potential error
	_ = r.ParseForm()

	emailStr := r.PostForm.Get("email")
	passwordStr := r.PostForm.Get("password")
	passwordConfirmStr := r.PostForm.Get("password-confirm")

	if emailStr == "" || passwordStr == "" || passwordConfirmStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = errors.BoxError("Please fill out all required fields").Render(r.Context(), w)
		return
	}

	email, err := valueobjects.NewEmail(emailStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errMsg := fmt.Sprintf("Invalid email: %s", err)
		_ = errors.BoxError(errMsg).Render(r.Context(), w)
		return

	}

	password, err := valueobjects.NewPassword(passwordStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errMsg := fmt.Sprintf("Password is not strong enough: %s", err)
		_ = errors.BoxError(errMsg).Render(r.Context(), w)
		return
	}

	if passwordStr != passwordConfirmStr {
		w.WriteHeader(http.StatusBadRequest)
		_ = errors.BoxError("Password and confirmation must match!").Render(r.Context(), w)
		return
	}

	uow := h.uowProvider.Provide()
	if err := uow.Begin(r.Context()); err != nil {
		l.Error("could not start transaction", slog.String("errMsg", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID, err := h.authService.Register(r.Context(), uow, email, password)
	if err == nil {
		// TODO: this should be handled in a transaction, ...
		// TODO: change this to free plan
		if err := h.permService.Create(r.Context(), uow, userID, permissions.PlanPro); err != nil {
			l.Error("failed to create user in perm service", slog.String("errMsg", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := uow.Commit(); err != nil {
			_ = uow.Rollback()
			l.Error("could not commit transaction", slog.String("errMsg", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// successfull register, redirect to /login (htmx)
		w.Header().Add("HX-Redirect", "/login")
		return
	}

	l.Error("failed to register", slog.String("errMsg", err.Error()))

	// handle error
	var errMsg string
	var statusCode int
	switch err {
	case auth.ErrUserAlreadyExists:
		errMsg = "This email is not available"
		statusCode = http.StatusConflict
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	_ = errors.BoxError(errMsg).Render(r.Context(), w)
}
