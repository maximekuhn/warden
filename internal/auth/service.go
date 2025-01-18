package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/valueobjects"
)

type AuthService struct {
	b Backend
}

func NewAuthService(b Backend) *AuthService {
	return &AuthService{b}
}

func (s *AuthService) Register(
	ctx context.Context,
	email valueobjects.Email,
	password valueobjects.Password,
) error {
	hashedPassword, err := bcryptHash(password)
	if err != nil {
		return err
	}
	user := NewUser(
		uuid.New(),
		email,
		hashedPassword,
		time.Now(),
		"",
		time.Unix(0, 0))
	return s.b.Save(ctx, *user)
}

func (s *AuthService) Login(
	ctx context.Context,
	email valueobjects.Email,
	password valueobjects.Password) (*http.Cookie, error) {
	user, found, err := s.b.GetByEmail(ctx, email)
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
		MaxAge:   int(time.Until(cookieExpiryDate)),
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
	if err := s.b.Update(ctx, *user, *updated); err != nil {
		return nil, err
	}
	return &cookie, nil
}

func (s *AuthService) Logout(
	ctx context.Context,
	cookie http.Cookie) error {
	if !validateCookieValue(cookie.Value) {
		return ErrCookieValueMalformed
	}
	user, found, err := s.b.GetBySessionId(ctx, cookie.Value)
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
	return s.b.Update(ctx, *user, *updated)
}

func (s *AuthService) Authenticate(
	ctx context.Context,
	cookie http.Cookie) (*User, error) {
	if !validateCookieValue(cookie.Value) {
		return nil, ErrCookieValueMalformed
	}
	user, found, err := s.b.GetBySessionId(ctx, cookie.Value)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ErrUserNotFound
	}
	if user.SessionId != cookie.Value {
		return nil, ErrBadSessionId
	}
	if user.SessionExpireDate.Unix() > time.Now().Unix() {
		return nil, ErrSessionExpired
	}
	return user, nil
}

func validateCookieValue(val string) bool {
	_, err := uuid.Parse(val)
	return err == nil
}
