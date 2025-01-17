package valueobjects

import "net/mail"

type Email struct {
	val string
}

func NewEmail(val string) (Email, error) {
	email := Email{}
	m, err := mail.ParseAddress(val)
	if err != nil {
		return email, err
	}
	email.val = m.Address
	return email, nil
}

func (e Email) Value() string {
	return e.val
}
