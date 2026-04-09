package pgtest

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type PgTestDatabase struct {
	pgTest   *pgTest
	released bool

	config pgTestConfig
	creds  PgTestCredentials

	adminPool *pgxpool.Pool
	pool      *pgxpool.Pool

	interDBMutex    *sync.Mutex
	releaseCallback func(db *PgTestDatabase) error
}

func newDatabase(pgTest *pgTest, config pgTestConfig, creds PgTestCredentials, releaseCallback func(db *PgTestDatabase) error) (*PgTestDatabase, error) {
	pool, err := pgxpool.New(context.Background(), creds.String())
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return &PgTestDatabase{
		pgTest:          pgTest,
		config:          config,
		creds:           creds,
		adminPool:       pgTest.pool,
		interDBMutex:    pgTest.interDBMutex,
		pool:            pool,
		released:        false,
		releaseCallback: releaseCallback,
	}, nil
}

func (d *PgTestDatabase) MigrateUp() error {
	if d.released {
		return ErrDatabaseReleased
	}

	migrPath := filepath.Join(d.config.WorkingDirectory, d.config.MigrationDirectory)

	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(d.pool)
	defer db.Close()

	err = goose.Up(db, migrPath)
	if err != nil {
		return err
	}

	return nil
}

func (d *PgTestDatabase) MigrateDown() error {
	if d.released {
		return ErrDatabaseReleased
	}

	migrPath := filepath.Join(d.config.WorkingDirectory, d.config.MigrationDirectory)

	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(d.pool)
	defer db.Close()

	err = goose.Down(db, migrPath)
	if err != nil {
		return err
	}

	return nil
}

func (d *PgTestDatabase) CreateTemplate() error {
	if d.released {
		return ErrNoTemplate
	}

	ctx := context.Background()

	if d.templateExists() {
		return ErrTemplateExists
	}

	_, _ = d.adminPool.Exec(ctx, `
		SELECT pg_terminate_backend(pid)
		FROM pg_stat_activity
		WHERE datname = $1
		`, d.creds.Database)

	_, err := d.adminPool.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s TEMPLATE %s", d.templateName(), d.creds.Database))
	if err != nil {
		return err
	}

	err = d.updatePool()
	if err != nil {
		return err
	}

	return nil
}

func (d *PgTestDatabase) ResetTemplate() error {
	if d.released {
		return ErrDatabaseReleased
	}
	ctx := context.Background()

	if !d.templateExists() {
		return ErrNoTemplate
	}

	_, _ = d.adminPool.Exec(ctx, `
		SELECT pg_terminate_backend(pid)
		FROM pg_stat_activity
		WHERE datname = $1
		`, d.templateName())

	_, err := d.adminPool.Exec(ctx, fmt.Sprintf("DROP DATABASE %s", d.templateName()))
	if err != nil {
		return err
	}

	return nil
}

func (d *PgTestDatabase) ApplyTemplate() error {
	if d.released {
		return ErrDatabaseReleased
	}
	ctx := context.Background()

	if !d.templateExists() {
		return ErrNoTemplate
	}

	d.pool.Close()
	d.pool = nil

	// Уничтожаем текущую базу
	_, _ = d.adminPool.Exec(ctx, `
		SELECT pg_terminate_backend(pid)
		FROM pg_stat_activity
		WHERE datname = $1
		`, d.creds.Database)

	_, err := d.adminPool.Exec(ctx, fmt.Sprintf("DROP DATABASE %s", d.creds.Database))
	if err != nil {
		return err
	}

	// Создаем базу из шаблона
	_, err = d.adminPool.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s WITH TEMPLATE %s", d.creds.Database, d.templateName()))
	if err != nil {
		return err
	}

	err = d.updatePool()
	if err != nil {
		return err
	}

	return nil
}

func (d *PgTestDatabase) templateName() string {
	return d.creds.Database + "_template"
}

func (d *PgTestDatabase) templateExists() bool {
	var exists bool

	err := d.adminPool.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT 1 FROM pg_database WHERE datname = $1
		)
	`, d.templateName()).Scan(&exists)

	if err != nil {
		return false
	}

	return exists
}

func (d *PgTestDatabase) updatePool() error {
	if d.pool != nil {
		d.pool.Close()
		d.pool = nil
	}

	pool, err := pgxpool.New(context.Background(), d.creds.String())
	if err != nil {
		return err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return err
	}

	d.pool = pool

	return nil
}

func (d *PgTestDatabase) Pool() *pgxpool.Pool {
	return d.pool
}

func (d *PgTestDatabase) Release() error {
	if d.released {
		return ErrDatabaseReleased
	}
	d.pool.Close()
	d.pool = nil

	err := d.pgTest.releaseDatabase(d)
	if err != nil {
		return err
	}

	d.released = true

	if d.releaseCallback != nil {
		return d.releaseCallback(d)
	}

	return nil
}

var (
	ErrDatabaseReleased = errors.New("database is released")
	ErrDatabaseNotFound = errors.New("database not found")
	ErrNoTemplate       = errors.New("no template")
	ErrTemplateExists   = errors.New("template exists")
)
