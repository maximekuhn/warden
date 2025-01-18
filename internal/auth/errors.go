package auth

import "errors"

// generic errors exposed by the auth package.
// These errors are mainly here for debug purposes, and should not be
// returned as-is to user. Otherwise, attackers could use them as a way
// to get more knowledge about the user base and potentially access
// an account.
var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyLoggedOut = errors.New("user already logged out")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrBadCredentials       = errors.New("bad credentials")
	ErrCookieValueMalformed = errors.New("cookie value malformed")
	ErrSessionExpired       = errors.New("session has expired")
	ErrBadSessionId         = errors.New("incorrect session id")
)
