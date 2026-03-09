package pgtx

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgTxConfig struct {
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

type PgTxOption func(*pgTxConfig)

func newConfig() pgTxConfig {
	return pgTxConfig{
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

// PgTxOption задаёт функциональную опцию конфигурации.
func (c *pgTxConfig) Option(opts ...PgTxOption) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithHost(host string) PgTxOption {
	return func(c *pgTxConfig) {
		c.host = host
	}
}

func WithPort(port uint16) PgTxOption {
	return func(c *pgTxConfig) {
		c.port = port
	}
}

func WithUser(user string) PgTxOption {
	return func(c *pgTxConfig) {
		c.user = user
	}
}

func WithPassword(password string) PgTxOption {
	return func(c *pgTxConfig) {
		c.password = password
	}
}

func WithDatabase(database string) PgTxOption {
	return func(c *pgTxConfig) {
		c.database = database
	}
}

func WithMaxConns(maxConns uint) PgTxOption {
	return func(c *pgTxConfig) {
		c.maxConns = maxConns
	}
}

// connLifeTime — примеры значений: 1s, 1m, 1h30m, 1d
func WithConnLifeTime(connLifeTime string) PgTxOption {
	return func(c *pgTxConfig) {
		c.connLifeTime = connLifeTime
	}
}

func WithRetryCount(retryCount uint) PgTxOption {
	return func(c *pgTxConfig) {
		c.retryCount = retryCount
	}
}

// New создаёт пул соединений PostgreSQL,
// DB-обёртку и менеджер транзакций.
func New(opts ...PgTxOption) (*DB, *PgTxManager, error) {
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

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, nil, err
	}

	return newDB(pool), newTxManager(pool, int(config.retryCount)), nil
}
