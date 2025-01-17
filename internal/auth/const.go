package auth

import "time"

const (
	CookieName = "warden-cookie"

	sessionValidityPeriod = 24 * time.Hour
)
