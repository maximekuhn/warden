package permissions

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserNotInServer = errors.New("user not in server")
)
