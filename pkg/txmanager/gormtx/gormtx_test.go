//go:build integration

package gormtx_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"media-equipment-tracker/pkg/pgtest"
	"media-equipment-tracker/pkg/txmanager/gormtx"
)

type GormTxTestSuite struct {
	suite.Suite

	db *gorm.DB
	pg *pgtest.PgTestDatabase
	tx *gormtx.GormTxManager
	dg *gormtx.DBGetter
}

func (s *GormTxTestSuite) SetupSuite() {
	pg, err := pgtest.Database()
	s.Require().NoError(err)

	s.pg = pg
	s.Require().NoError(pg.MigrateUp())

	db := stdlib.OpenDBFromPool(s.pg.Pool())
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	s.Require().NoError(err)

	s.Require().NoError(gdb.Exec("CREATE TABLE IF NOT EXISTS test_table (id SERIAL PRIMARY KEY, name TEXT)").Error)

	s.Require().NoError(pg.CreateTemplate())
}

func (s *GormTxTestSuite) TearDownSuite() {
	_ = s.pg.Release()
}

func (s *GormTxTestSuite) SetupTest() {
	s.Require().NoError(s.pg.ApplyTemplate())

	db := stdlib.OpenDBFromPool(s.pg.Pool())
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	s.Require().NoError(err)

	s.db = gdb
	s.dg = gormtx.NewDBGetter(gdb)
	s.tx = gormtx.NewGormTxManager(gdb, 3) // retry count 3
}

// Test 1: Simple transaction in dbgetter
func (s *GormTxTestSuite) TestSimpleTransactionDBGetter() {
	tx := s.db.Begin()
	defer tx.Rollback()

	// Simulate insert
	err := tx.Exec("INSERT INTO test_table (name) VALUES (?)", "test").Error
	s.Require().NoError(err)

	err = tx.Commit().Error
	s.NoError(err)

	// Check if inserted
	var count int
	err = s.db.Raw("SELECT COUNT(*) FROM test_table WHERE name = ?", "test").Scan(&count).Error
	s.NoError(err)
	s.Equal(1, count)
}

// Test 2: Failed transaction in dbgetter
func (s *GormTxTestSuite) TestFailedTransactionDBGetter() {
	tx := s.db.Begin()
	defer tx.Rollback()

	// Simulate insert
	err := tx.Exec("INSERT INTO test_table (name) VALUES (?)", "failed").Error
	s.Require().NoError(err)

	// Rollback
	tx.Rollback()

	// Check not inserted
	var count int
	err = s.db.Raw("SELECT COUNT(*) FROM test_table WHERE name = ?", "failed").Scan(&count).Error
	s.NoError(err)
	s.Equal(0, count)
}

// Test 3: Simple transaction withinTX
func (s *GormTxTestSuite) TestSimpleTransactionWithinTX() {
	ctx := context.Background()

	err := s.tx.WithinTx(ctx, func(ctx context.Context) error {
		db, err := s.dg.GetDB(ctx)
		if err != nil {
			return err
		}
		return db.Exec("INSERT INTO test_table (name) VALUES (?)", "within").Error
	})
	s.NoError(err)

	// Check inserted
	var count int
	err = s.db.Raw("SELECT COUNT(*) FROM test_table WHERE name = ?", "within").Scan(&count).Error
	s.NoError(err)
	s.Equal(1, count)
}

// Test 4: Failed transaction withinTX
func (s *GormTxTestSuite) TestFailedTransactionWithinTX() {
	ctx := context.Background()

	err := s.tx.WithinTx(ctx, func(ctx context.Context) error {
		db, err := s.dg.GetDB(ctx)
		if err != nil {
			return err
		}
		_ = db.Exec("INSERT INTO test_table (name) VALUES (?)", "failed_within")
		return errors.New("forced error")
	})
	s.Error(err)

	// Check not inserted
	var count int
	err = s.db.Raw("SELECT COUNT(*) FROM test_table WHERE name = ?", "failed_within").Scan(&count).Error
	s.NoError(err)
	s.Equal(0, count)
}

// Test 5: Nested transaction dbgetter
func (s *GormTxTestSuite) TestNestedTransactionDBGetter() {
	ctx := context.Background()

	db, err := s.dg.GetDB(ctx)
	s.Require().NoError(err)

	tx := db.Begin()
	defer tx.Rollback()

	// Outer insert
	err = tx.Exec("INSERT INTO test_table (name) VALUES (?)", "outer").Error
	s.Require().NoError(err)

	// Inner transaction
	tx.Transaction(func(tx *gorm.DB) error {
		err = tx.Exec("INSERT INTO test_table (name) VALUES (?)", "inner").Error
		s.Require().NoError(err)
		return nil
	})

	err = tx.Commit().Error
	s.NoError(err)

	// Check both inserted
	var countOuter, countInner int
	_ = s.db.Raw("SELECT COUNT(*) FROM test_table WHERE name = ?", "outer").Scan(&countOuter).Error
	_ = s.db.Raw("SELECT COUNT(*) FROM test_table WHERE name = ?", "inner").Scan(&countInner).Error
	s.Equal(1, countOuter)
	s.Equal(1, countInner)
}

// Test 6: Nested transaction dbgetter, one fails but others pass
func (s *GormTxTestSuite) TestNestedTransactionDBGetterOneFailsOthersPass() {
	ctx := context.Background()

	db, err := s.dg.GetDB(ctx)
	s.Require().NoError(err)

	tx := db.Begin()
	defer tx.Rollback()

	// Outer insert
	err = tx.Exec("INSERT INTO test_table (name) VALUES (?)", "outer2").Error
	s.Require().NoError(err)

	tx.Transaction(func(tx *gorm.DB) error {
		err = tx.Exec("INSERT INTO test_table (name) VALUES (?)", "inner_fail").Error
		s.Require().NoError(err)
		return errors.New("propagated error")
	})

	tx.Transaction(func(tx *gorm.DB) error {
		err = tx.Exec("INSERT INTO test_table (name) VALUES (?)", "inner").Error
		s.Require().NoError(err)
		return nil
	})

	// Commit outer
	err = tx.Commit().Error
	s.NoError(err)

	// Check outer inserted, inner not
	var countOuter, countInnerFail, countInner int
	_ = s.db.Raw("SELECT COUNT(*) FROM test_table WHERE name = ?", "outer2").Scan(&countOuter).Error
	_ = s.db.Raw("SELECT COUNT(*) FROM test_table WHERE name = ?", "inner_fail").Scan(&countInnerFail).Error
	_ = s.db.Raw("SELECT COUNT(*) FROM test_table WHERE name = ?", "inner").Scan(&countInner).Error
	s.Equal(1, countOuter)
	s.Equal(0, countInnerFail)
	s.Equal(1, countInner)

}

// Test 7: Nested transaction dbgetter, one fails, others don't commit
func (s *GormTxTestSuite) TestNestedTransactionDBGetterOneFailsCancelsOuter() {
	ctx := context.Background()

	db, err := s.dg.GetDB(ctx)
	s.Require().NoError(err)

	tx := db.Begin()
	defer tx.Rollback()

	// Outer insert
	err = tx.Exec("INSERT INTO test_table (name) VALUES (?)", "outer").Error
	s.Require().NoError(err)

	// Inner transaction that fails and propagates
	err = tx.Transaction(func(tx *gorm.DB) error {
		err = tx.Exec("INSERT INTO test_table (name) VALUES (?)", "inner_fail").Error
		s.Require().NoError(err)
		return errors.New("propagated error")
	})
	s.Error(err)

	err = tx.Rollback().Error
	s.NoError(err)

	// Check nothing inserted
	var countOuter, countInner int
	_ = s.db.Raw("SELECT COUNT(*) FROM test_table WHERE name = ?", "outer").Scan(&countOuter).Error
	_ = s.db.Raw("SELECT COUNT(*) FROM test_table WHERE name = ?", "inner_fail").Scan(&countInner).Error
	s.Equal(0, countOuter)
	s.Equal(0, countInner)
}

// Test 8: Nested transaction with withinTx outer
func (s *GormTxTestSuite) TestNestedTransactionWithinTxOuter() {
	ctx := context.Background()

	err := s.tx.WithinTx(ctx, func(ctx context.Context) error {
		db, err := s.dg.GetDB(ctx)
		if err != nil {
			return err
		}
		err = db.Exec("INSERT INTO test_table (name) VALUES (?)", "outer_within").Error
		s.Require().NoError(err)

		err = db.Transaction(func(tx *gorm.DB) error {
			err = tx.Exec("INSERT INTO test_table (name) VALUES (?)", "inner_within").Error
			s.Require().NoError(err)
			return nil
		})
		s.Require().NoError(err)
		return nil
	})
	s.NoError(err)

	// Check both inserted
	var countOuter, countInner int
	_ = s.db.Raw("SELECT COUNT(*) FROM test_table WHERE name = ?", "outer_within").Scan(&countOuter).Error
	_ = s.db.Raw("SELECT COUNT(*) FROM test_table WHERE name = ?", "inner_within").Scan(&countInner).Error
	s.Equal(1, countOuter)
	s.Equal(1, countInner)
}

// Test 9: Nested transaction with withinTx outer, inner fails but outer passes
func (s *GormTxTestSuite) TestNestedTransactionWithinTxOneFailsOthersPass() {
	ctx := context.Background()

	err := s.tx.WithinTx(ctx, func(ctx context.Context) error {
		db, err := s.dg.GetDB(ctx)
		if err != nil {
			return err
		}
		err = db.Exec("INSERT INTO test_table (name) VALUES (?)", "outer_within").Error
		s.Require().NoError(err)

		// Inner transaction that fails but we handle it
		err = db.Transaction(func(tx *gorm.DB) error {
			err = tx.Exec("INSERT INTO test_table (name) VALUES (?)", "inner_fail_within").Error
			s.Require().NoError(err)
			return errors.New("inner forced error")
		})
		s.Require().Error(err)

		err = db.Transaction(func(tx *gorm.DB) error {
			err = tx.Exec("INSERT INTO test_table (name) VALUES (?)", "inner_within").Error
			s.Require().NoError(err)
			return nil
		})
		// Don't return innerErr, so outer commits
		return nil
	})
	s.NoError(err) // Outer succeeds

	// Check outer inserted, inner not
	var countOuter, countInnerFail, countInner int
	_ = s.db.Raw("SELECT COUNT(*) FROM test_table WHERE name = ?", "outer_within").Scan(&countOuter).Error
	_ = s.db.Raw("SELECT COUNT(*) FROM test_table WHERE name = ?", "inner_fail_within").Scan(&countInnerFail).Error
	_ = s.db.Raw("SELECT COUNT(*) FROM test_table WHERE name = ?", "inner_within").Scan(&countInner).Error
	s.Equal(1, countOuter)
	s.Equal(0, countInnerFail)
	s.Equal(1, countInner)
}

// Test 10: Nested transaction with withinTx outer, inner fails, outer cancels
func (s *GormTxTestSuite) TestNestedTransactionWithinTxOneFailsCancelsOuter() {
	ctx := context.Background()

	err := s.tx.WithinTx(ctx, func(ctx context.Context) error {
		db, err := s.dg.GetDB(ctx)
		if err != nil {
			return err
		}
		err = db.Exec("INSERT INTO test_table (name) VALUES (?)", "outer_within").Error

		// Inner transaction that succeeds
		err = db.Transaction(func(tx *gorm.DB) error {
			err = tx.Exec("INSERT INTO test_table (name) VALUES (?)", "inne_within").Error
			s.Require().NoError(err)
			return nil
		})

		// Inner transaction that fails and propagates
		err = db.Transaction(func(tx *gorm.DB) error {
			err = tx.Exec("INSERT INTO test_table (name) VALUES (?)", "inner_fail_within").Error
			s.Require().NoError(err)
			return errors.New("inner forced error")
		})
		s.Require().Error(err)
		return err
	})
	s.Error(err) // Outer fails

	// Check nothing inserted
	var countOuter, countInnerFail, countInner int
	_ = s.db.Raw("SELECT COUNT(*) FROM test_table WHERE name = ?", "outer_within").Scan(&countOuter).Error
	_ = s.db.Raw("SELECT COUNT(*) FROM test_table WHERE name = ?", "inner_fail_within").Scan(&countInner).Error
	_ = s.db.Raw("SELECT COUNT(*) FROM test_table WHERE name = ?", "inner_within").Scan(&countInner).Error
	s.Equal(0, countOuter)
	s.Equal(0, countInner)
	s.Equal(0, countInnerFail)
}

func TestGormTxTestSuite(t *testing.T) {
	suite.Run(t, new(GormTxTestSuite))
}
