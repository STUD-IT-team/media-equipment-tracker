//go:build integration

package pgtest_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"media-equipment-tracker/pkg/pgtest"
)

func TestPgTest_BasicLifecycle(t *testing.T) {
	db, err := pgtest.Database()
	require.NoError(t, err)
	require.NotNil(t, db)

	defer func() {
		require.NoError(t, db.Release())
	}()

	ctx := context.Background()

	// проверяем что pool работает
	var one int
	err = db.Pool().QueryRow(ctx, "SELECT 1").Scan(&one)
	require.NoError(t, err)
	require.Equal(t, 1, one)
}

func TestPgTest_Migrations(t *testing.T) {
	db, err := pgtest.Database()
	require.NoError(t, err)
	defer db.Release()

	err = db.MigrateUp()
	require.NoError(t, err)
}

func TestPgTest_TemplateFlow(t *testing.T) {
	db, err := pgtest.Database()
	require.NoError(t, err)
	defer db.Release()

	ctx := context.Background()

	// создаём тестовую таблицу
	_, err = db.Pool().Exec(ctx, `CREATE TABLE test_table (id INT)`)
	require.NoError(t, err)

	// вставляем данные
	_, err = db.Pool().Exec(ctx, `INSERT INTO test_table (id) VALUES (42)`)
	require.NoError(t, err)

	// создаём template
	require.NoError(t, db.CreateTemplate())

	// мутируем данные
	_, err = db.Pool().Exec(ctx, `DELETE FROM test_table`)
	require.NoError(t, err)

	// применяем template (должно вернуть данные)
	require.NoError(t, db.ApplyTemplate())

	var count int
	err = db.Pool().QueryRow(ctx, `SELECT COUNT(*) FROM test_table`).Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 1, count)
}

func TestPgTest_MultipleDatabases(t *testing.T) {
	db1, err := pgtest.Database()
	require.NoError(t, err)
	defer db1.Release()

	db2, err := pgtest.Database()
	require.NoError(t, err)
	defer db2.Release()

	ctx := context.Background()

	// создаём таблицу только в db1
	_, err = db1.Pool().Exec(ctx, `CREATE TABLE only_here (id INT)`)
	require.NoError(t, err)

	// в db2 её быть не должно
	_, err = db2.Pool().Exec(ctx, `SELECT * FROM only_here`)
	require.Error(t, err)
}
