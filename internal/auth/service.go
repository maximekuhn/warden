package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
)

type AuthService struct {
	b Backend
}

func NewAuthService(b Backend) *AuthService {
	return &AuthService{b}
}

func (s *AuthService) Register(
	ctx context.Context,
	uow transaction.UnitOfWork,
	email valueobjects.Email,
	password valueobjects.Password,
) (uuid.UUID, error) {
	hashedPassword, err := bcryptHash(password)
	if err != nil {
		return uuid.UUID{}, err
	}
	userID := uuid.New()
	user := NewUser(
		userID,
		email,
		hashedPassword,
		time.Now(),
		"",
		time.Unix(0, 0))
	if err := s.b.Save(ctx, uow, *user); err != nil {
		return uuid.UUID{}, err
	}
	return userID, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	uow transaction.UnitOfWork,
	email valueobjects.Email,
	password valueobjects.Password) (*http.Cookie, error) {
	user, found, err := s.b.GetByEmail(ctx, uow, email)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ErrUserNotFound
	}
	if !bcryptVerify(password, user.HashedPassord) {
		return nil, ErrBadCredentials
	}

	// create new session and associated cookie
	sessionId := uuid.NewString()
	now := time.Now()
	cookieExpiryDate := now.Add(sessionValidityPeriod)
	cookie := http.Cookie{
		Name:     CookieName,
		Value:    sessionId,
		MaxAge:   int(time.Until(cookieExpiryDate).Seconds()),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	// update user information
	updated := NewUser(
		user.ID,
		email,
		user.HashedPassord,
		user.CreatedAt,
		sessionId,
		cookieExpiryDate)
	if err := s.b.Update(ctx, uow, *user, *updated); err != nil {
		return nil, err
	}
	return &cookie, nil
}

func (s *AuthService) Logout(
	ctx context.Context,
	uow transaction.UnitOfWork,
	cookie http.Cookie) error {
	if !validateCookieValue(cookie.Value) {
		return ErrCookieValueMalformed
	}
	user, found, err := s.b.GetBySessionId(ctx, uow, cookie.Value)
	if err != nil {
		return err
	}
	if !found {
		return ErrUserNotFound
	}
	if user.IsLoggedOut() {
		return ErrUserAlreadyLoggedOut
	}

	updated := NewUser(
		user.ID,
		user.Email,
		user.HashedPassord,
		user.CreatedAt,
		"",
		time.Unix(0, 0))
	return s.b.Update(ctx, uow, *user, *updated)
}

func (s *AuthService) Authenticate(
	ctx context.Context,
	uow transaction.UnitOfWork,
	cookie http.Cookie) (*User, error) {
	if !validateCookieValue(cookie.Value) {
		return nil, ErrCookieValueMalformed
	}
	user, found, err := s.b.GetBySessionId(ctx, uow, cookie.Value)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ErrUserNotFound
	}
	if user.SessionId != cookie.Value {
		return nil, ErrBadSessionId
	}
	if user.SessionExpireDate.Before(time.Now()) {
		return nil, ErrSessionExpired
	}
	return user, nil
}

func validateCookieValue(val string) bool {
	_, err := uuid.Parse(val)
	return err == nil
}
