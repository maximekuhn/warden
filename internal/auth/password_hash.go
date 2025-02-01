package auth

import (
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
	"golang.org/x/crypto/bcrypt"
)

type HashedPassword []byte

func bcryptHash(p valueobjects.Password) (HashedPassword, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(p.Value()), bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, err
	}
	return HashedPassword(hash), nil
}

func bcryptVerify(p valueobjects.Password, h HashedPassword) bool {
	return bcrypt.CompareHashAndPassword(h, []byte(p.Value())) == nil
}
