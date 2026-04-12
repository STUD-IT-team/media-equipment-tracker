package authservice

type TokenRepository interface {
	Check(token string) bool
	Add(token string) bool
	Delete(token string)
}
