package pgtest

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type pgTest struct {
	container testcontainers.Container
	config    *pgTestConfig

	creds PgTestCredentials
	pool  *pgxpool.Pool

	mutex   *sync.Mutex
	started bool

	databases    []*PgTestDatabase
	interDBMutex *sync.Mutex
	refCount     int
}

func newPgTest(config *pgTestConfig) *pgTest {
	return &pgTest{
		config:       config,
		mutex:        &sync.Mutex{},
		started:      false,
		databases:    make([]*PgTestDatabase, 0),
		interDBMutex: &sync.Mutex{},
	}
}

const (
	pgTestPort = "5432/tcp"
)

func (p *pgTest) Start() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	ctx := context.Background()

	fmt.Println(p.config)

	if p.started {
		return nil
	}

	var err error
	if p.config.External {
		err = p.startExternal()
	} else {
		err = p.startContainer()
	}

	if err != nil {
		return err
	}

	p.pool, err = pgxpool.New(ctx, p.creds.String())
	if err != nil {
		return err
	}

	err = p.pool.Ping(ctx)
	if err != nil {
		return err
	}

	p.started = true

	return nil
}

func (p *pgTest) startContainer() error {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        p.config.Image,
		ExposedPorts: []string{pgTestPort},
		Env: map[string]string{
			"POSTGRES_USER":     p.config.Username,
			"POSTGRES_PASSWORD": p.config.Password,
			"POSTGRES_DB":       p.config.Database,
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort(nat.Port(pgTestPort)),
		),
	}

	cnt, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return err
	}

	host, err := cnt.Host(ctx)
	if err != nil {
		return err
	}
	port, err := cnt.MappedPort(ctx, nat.Port(pgTestPort))
	if err != nil {
		return err
	}

	p.creds = fromConfig(p.config, uint16(port.Int()), host)
	p.container = cnt

	return nil
}

func (p *pgTest) startExternal() error {
	p.creds = fromConfig(p.config, p.config.Port, p.config.Host)
	return nil
}

func (p *pgTest) Stop() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.started {
		return ErrPgTestNotStarted
	}

	if p.refCount > 0 {
		return ErrPgTestInUse
	}

	p.pool.Close()
	p.pool = nil

	var err error
	if p.config.External {
		err = p.stopExternal()
	} else {
		err = p.stopContainer()
	}

	if err != nil {
		return err
	}

	p.started = false

	return nil
}

func (p *pgTest) stopContainer() error {
	err := p.container.Terminate(context.Background())
	if err != nil {
		return err
	}

	p.container = nil
	return nil
}

func (p *pgTest) stopExternal() error {
	return nil
}

func (p *pgTest) Database() (*PgTestDatabase, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.started {
		return nil, ErrPgTestNotStarted
	}

	dbName := "test_" + strings.ReplaceAll(uuid.New().String(), "-", "_")

	err := p.createDatabase(dbName)
	if err != nil {
		return nil, err
	}

	creds := PgTestCredentials{
		Username: p.creds.Username,
		Password: p.creds.Password,
		Database: dbName,
		Port:     p.creds.Port,
		Host:     p.creds.Host,
	}

	db, err := newDatabase(p, p.config, creds, nil)
	if err != nil {
		return nil, err
	}

	p.databases = append(p.databases, db)
	p.refCount++

	return db, nil
}

func (p *pgTest) DatabaseWithCallback(callback func(db *PgTestDatabase) error) (*PgTestDatabase, error) {
	db, err := p.Database()
	if err != nil {
		return nil, err
	}

	db.releaseCallback = callback

	return db, nil
}

func (p *pgTest) releaseDatabase(db *PgTestDatabase) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	found := false
	for i, d := range p.databases {
		if d == db {
			p.databases = append(p.databases[:i], p.databases[i+1:]...)
			p.refCount--
			found = true
			break
		}
	}
	if !found {
		return ErrPgTestNotFound
	}

	return p.dropDatabase(db.creds.Database)
}

func (p *pgTest) createDatabase(name string) error {
	_, err := p.pool.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", name))
	if err != nil {
		return err
	}

	return nil
}

func (p *pgTest) dropDatabase(name string) error {
	_, err := p.pool.Exec(context.Background(), fmt.Sprintf("DROP DATABASE %s", name))
	if err != nil {
		return err
	}

	return nil
}

var (
	ErrPgTestNotStarted = errors.New("pgtest not started")
	ErrPgTestInUse      = errors.New("pgtest is still in use")
	ErrPgTestNotFound   = errors.New("pgtest not found")
)
