package valueobjects

import (
	"errors"
	"fmt"
	"slices"
	"unicode"
)

type Password struct {
	val string
}

func NewPassword(val string) (Password, error) {
	// XXX: validation could be stronger
	const (
		minLength = 12
		maxLength = 72
	)
	var symbols = []rune{'+', '!', '@', '#', '^', '=', '-', '_'}
	p := Password{}
	if len(val) < minLength {
		return p, fmt.Errorf("password must be at least %d characters", minLength)
	}
	if len(val) > maxLength {
		return p, fmt.Errorf("password must be at most %d characters", minLength)
	}
	if !containsUppercase(val) {
		return p, errors.New("password must contain at least one uppercase letter")
	}
	if !containsNumber(val) {
		return p, errors.New("password must contain at least one number")
	}
	if !containsSpecial(val, symbols) {
		return p, fmt.Errorf("password must contain at least one special characters (allowed: %v)", symbols)
	}
	p.val = val
	return p, nil
}

func (p Password) Value() string {
	return p.val
}

func containsUppercase(val string) bool {
	for _, c := range val {
		if unicode.IsUpper(c) {
			return true
		}
	}
	return false
}

func containsNumber(val string) bool {
	for _, c := range val {
		if unicode.IsNumber(c) {
			return true
		}
	}
	return false
}

func containsSpecial(val string, symbols []rune) bool {
	for _, c := range val {
		if slices.Contains(symbols, c) {
			return true
		}
	}
	return false
}
