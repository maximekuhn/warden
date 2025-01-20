package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/maximekuhn/warden/internal/auth"
	"github.com/maximekuhn/warden/internal/commands"
	"github.com/maximekuhn/warden/internal/handlers/handlerutils"
	"github.com/maximekuhn/warden/internal/logger"
	"github.com/maximekuhn/warden/internal/middlewares"
	uierrors "github.com/maximekuhn/warden/internal/ui/components/errors"
	"github.com/maximekuhn/warden/internal/ui/pages"
)

type SignupHandler struct {
	logger               *slog.Logger
	createUserCmdHandler *commands.CreateUserCommandHandler
}

func NewSignupHandler(
	l *slog.Logger,
	cmdHandler *commands.CreateUserCommandHandler,
) *SignupHandler {
	return &SignupHandler{
		logger:               l,
		createUserCmdHandler: cmdHandler,
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

	// parse form values
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	emailStr := r.PostForm.Get("email")
	passwordStr := r.PostForm.Get("password")
	passwordConfirmStr := r.PostForm.Get("password-confirm")
	email, err := handlerutils.ToEmailOrReturnErrorBox(w, emailStr)
	if err != nil {
		l.Info("invalid email", slog.String("reason", err.Error()))
		return
	}
	password, err := handlerutils.ToPasswordOrReturnErrorBox(w, passwordStr)
	if err != nil {
		l.Info("invalid password", slog.String("reason", err.Error()))
		return
	}
	passwordConfirm, err := handlerutils.ToPasswordOrReturnErrorBox(w, passwordConfirmStr)
	if err != nil {
		l.Info("invalid password confirmation", slog.String("reason", err.Error()))
		return
	}

	// handle user creation
	err = h.createUserCmdHandler.Handle(r.Context(), commands.CreateUserCommand{
		Email:           email,
		Password:        password,
		PasswordConfirm: passwordConfirm,
	})
	if err == nil {
		l.Info("successfully created user")
		w.Header().Add("HX-Redirect", "/login")
		w.WriteHeader(http.StatusCreated)
		return
	}

	// oh no.. it failed ;(
	l.Error("failed to register", slog.String("errMsg", err.Error()))
	h.postHandleError(w, r, err)
}

func (h *SignupHandler) postHandleError(w http.ResponseWriter, r *http.Request, err error) {
	var statusCode int
	var errMsg string

	if errors.Is(err, commands.ErrPasswordAndConfirmationDontMatch) {
		statusCode = http.StatusBadRequest
		errMsg = "Password and confirmation must match"
	} else if errors.Is(err, auth.ErrUserAlreadyExists) {
		statusCode = http.StatusConflict
		errMsg = "This email is unavailable. Try another one"
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// return error box
	w.WriteHeader(statusCode)
	if err := uierrors.BoxError(errMsg).Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
