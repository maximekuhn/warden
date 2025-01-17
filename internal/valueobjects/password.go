package valueobjects

type Password struct {
	val string
}

func NewPassword(val string) (Password, error) {
	p := Password{}
	// TODO: validation
	p.val = val
	return p, nil
}

func (p Password) Value() string {
	return p.val
}
