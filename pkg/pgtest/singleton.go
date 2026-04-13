package pgtest

import (
	"errors"
	"sync"
)

var (
	instance *pgTest
	once     sync.Once
	initErr  error

	globalMu sync.Mutex
)

func getInstance() (*pgTest, error) {
	once.Do(func() {
		cfg := configFromEnv()
		instance = newPgTest(&cfg)
		initErr = instance.Start()
	})

	return instance, initErr
}

func Database() (*PgTestDatabase, error) {
	pg, err := getInstance()
	if err != nil {
		return nil, err
	}

	db, err := pg.DatabaseWithCallback(tryStop)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func tryStop(_ *PgTestDatabase) error {
	globalMu.Lock()
	defer globalMu.Unlock()

	if instance == nil {
		return nil
	}

	err := instance.Stop()
	if errors.Is(err, ErrPgTestInUse) {
		return nil
	} else if err != nil {
		return err
	}

	// сбрасываем singleton
	instance = nil
	once = sync.Once{}

	return nil
}
