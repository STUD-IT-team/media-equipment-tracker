package pgtest

import (
	"fmt"
)

type PgTestCredentials struct {
	Username string
	Password string
	Database string
	Port     uint16
	Host     string
}

func fromConfig(config *pgTestConfig, port uint16, host string) PgTestCredentials {
	return PgTestCredentials{
		Username: config.Username,
		Password: config.Password,
		Database: config.Database,
		Port:     port,
		Host:     host,
	}
}

func (c PgTestCredentials) String() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", c.Username, c.Password, c.Host, c.Port, c.Database)
}
