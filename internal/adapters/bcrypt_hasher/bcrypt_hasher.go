package bcrypthasher

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrHash          = errors.New("failed to hash password")
	ErrPassword      = errors.New("error CompareHashAndPassword")
	ErrEmptyPassword = errors.New("error the password cannot be empty")
)

type BcryptHasher struct {
}

func (h *BcryptHasher) HashPassword(password string) (string, error) {
	if password == "" {
		return "", ErrEmptyPassword
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrHash
	}
	return string(hashedPassword), nil
}

func (h *BcryptHasher) CheckPassword(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		err = ErrPassword
	}
	return err
}
