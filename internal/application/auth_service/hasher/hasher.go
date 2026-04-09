package hasher

import (
	bcrypthasher "media-equipment-tracker/internal/adapters/bcrypt_hasher"
)

type Hasher interface {
	HashPassword(password string) (string, error) // HashPassword возвращает bcrypt хэш пароля
	CheckPassword(password string, hashedPassword string) error
}

func NewHasher() (Hasher, error) {
	return &bcrypthasher.BcryptHasher{}, nil
}
