package auth

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyLoggedOut = errors.New("user already logged out")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrBadCredentials       = errors.New("bad credentials")
	ErrCookieValueMalformed = errors.New("cookie value malformed")
)
