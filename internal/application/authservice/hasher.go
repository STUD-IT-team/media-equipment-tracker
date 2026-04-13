package authservice

type Hasher interface {
	HashPassword(password string) (string, error) // HashPassword возвращает bcrypt хэш пароля
	CheckPassword(password string, hashedPassword string) error
}
