package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/maximekuhn/warden/internal/auth"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
	"github.com/maximekuhn/warden/internal/logger"
	"github.com/maximekuhn/warden/internal/middlewares"
	"github.com/maximekuhn/warden/internal/ui/components/errors"
	"github.com/maximekuhn/warden/internal/ui/pages"
)

type LoginHandler struct {
	logger      *slog.Logger
	service     *auth.AuthService
	uowProvider transaction.UnitOfWorkProvider
}

func NewLoginHandler(
	l *slog.Logger,
	s *auth.AuthService,
	uowProvider transaction.UnitOfWorkProvider,
) *LoginHandler {
	return &LoginHandler{
		logger:      l,
		service:     s,
		uowProvider: uowProvider,
	}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (h *LoginHandler) get(w http.ResponseWriter, r *http.Request) {
	if err := pages.Login().Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *LoginHandler) post(w http.ResponseWriter, r *http.Request) {
	l := logger.UpgradeLoggerWithRequestId(r.Context(), middlewares.RequestIdKey, h.logger)

	// TODO: handle potentiel error
	_ = r.ParseForm()

	emailStr := r.PostForm.Get("email")
	passwordStr := r.PostForm.Get("password")

	if emailStr == "" || passwordStr == "" {
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

	uow := h.uowProvider.Provide()

	cookie, err := h.service.Login(r.Context(), uow, email, password)
	if err != nil {
		h.handleLoginError(w, r, l, err)
		return
	}

	l.Info("login successfull")

	http.SetCookie(w, cookie)
	w.Header().Add("HX-Redirect", "/")
}

func (h *LoginHandler) handleLoginError(w http.ResponseWriter, r *http.Request, l *slog.Logger, err error) {
	l.Error("failed to login", slog.String("errMsg", err.Error()))

	var errMsg string
	var statusCode int
	switch err {
	case auth.ErrBadCredentials:
		fallthrough
	case auth.ErrUserNotFound:
		errMsg = "Bad credentials"
		statusCode = http.StatusBadRequest
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	fmt.Printf("statusCode: %v\n", statusCode)
	w.WriteHeader(statusCode)
	_ = errors.BoxError(errMsg).Render(r.Context(), w)
}
