package config

import "time"

const (
	PostgresHost            = "localhost"
	PostgresPort            = 5432
	PostgresUser            = "uuser"
	PostgresPassword        = "ppassword"
	PostgresDatabase        = "eqtracker"
	TokenSymmetricKey       = "12345678901234567890123456789012"
	AccessTokenDuration     = 24 * time.Hour
	APIVersion              = "/api/v1"
	AppPort                 = 8080
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	EnableFileLogging       = true
)
