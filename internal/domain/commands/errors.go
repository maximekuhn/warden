package commands

import "errors"

// generic errors returned by various commands.
// HTTP handlers can use them to return the correct status code to the caller.
var (
	ErrPasswordAndConfirmationDontMatch = errors.New("password and confirmation don't match")
)
