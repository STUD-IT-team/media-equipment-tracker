package gormtx

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type gormTxConfig struct {
	host     string
	port     uint16
	user     string
	password string
	database string

	maxConns     uint
	connLifeTime string
	tls          bool

	retryCount uint
}

type GormTxOption func(*gormTxConfig)

func newConfig() gormTxConfig {
	return gormTxConfig{
		host:     "localhost",
		port:     5432,
		user:     "postgres",
		password: "postgres",
		database: "postgres",

		maxConns:     10,
		connLifeTime: "5m",
		tls:          false,

		retryCount: 3,
	}
}

func (c *gormTxConfig) Option(opts ...GormTxOption) {
	for _, opt := range opts {
		opt(c)
	}
}

// GormTxOption задаёт функциональную опцию конфигурации.
func WithHost(host string) GormTxOption {
	return func(c *gormTxConfig) {
		c.host = host
	}
}

func WithPort(port uint16) GormTxOption {
	return func(c *gormTxConfig) {
		c.port = port
	}
}

func WithUser(user string) GormTxOption {
	return func(c *gormTxConfig) {
		c.user = user
	}
}

func WithPassword(password string) GormTxOption {
	return func(c *gormTxConfig) {
		c.password = password
	}
}

func WithDatabase(database string) GormTxOption {
	return func(c *gormTxConfig) {
		c.database = database
	}
}

func WithMaxConns(maxConns uint) GormTxOption {
	return func(c *gormTxConfig) {
		c.maxConns = maxConns
	}
}

// connLifeTime — примеры значений: 1s, 1m, 1h30m, 1d
func WithConnLifeTime(connLifeTime string) GormTxOption {
	return func(c *gormTxConfig) {
		c.connLifeTime = connLifeTime
	}
}

func WithRetryCount(retryCount uint) GormTxOption {
	return func(c *gormTxConfig) {
		c.retryCount = retryCount
	}
}

type GormTxOptionFunc func(*gormTxConfig)

func New(opts ...GormTxOption) (*DBGetter, *GormTxManager, error) {
	config := newConfig()
	config.Option(opts...)

	var tls string
	if config.tls {
		tls = "require"
	} else {
		tls = "disable"
	}

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s pool_max_conns=%d pool_max_conn_lifetime=%s",
		config.host, config.port, config.user, config.password, config.database, tls, config.maxConns, config.connLifeTime)

	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	return NewDBGetter(db), NewGormTxManager(db, int(config.retryCount)), nil
}
