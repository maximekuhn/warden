package valueobjects

import (
	"errors"
	"fmt"
	"strings"
)

type MinecraftServerName struct {
	val string
}

func NewMinecraftServerName(val string) (MinecraftServerName, error) {
	const (
		maxLength = 24
		minLength = 8
	)

	name := MinecraftServerName{}
	trimmedVal := strings.Trim(val, " ")
	if trimmedVal == "" {
		return name, errors.New("name can not be empty")
	}
	if len(trimmedVal) > maxLength {
		return name, fmt.Errorf("name can not exceed %d characters", maxLength)
	}
	if len(trimmedVal) < minLength {
		return name, fmt.Errorf("name must be at least %d characters", minLength)
	}
	name.val = trimmedVal
	return name, nil
}

func (m MinecraftServerName) Value() string {
	return m.val
}
