package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/kriuchkov/tonbeacon/pkg/containers"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type RepositoryTestSuite struct {
	suite.Suite
	db      *bun.DB
	adapter *DatabaseAdapter
	tables  []string
}

func (suite *RepositoryTestSuite) SetupSuite() {
	suite.db = setupTestDB(suite.T())
	suite.tables = getTables(suite.T(), suite.db)
	suite.adapter = New(suite.db)

}

func (suite *RepositoryTestSuite) TeardownSuite() {
	ctx := context.Background()
	for _, table := range suite.tables {
		_, err := suite.db.ExecContext(ctx, "TRUNCATE TABLE "+table+" CASCADE")
		require.NoError(suite.T(), err)
	}
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func setupTestDB(t *testing.T) *bun.DB {
	ctx := context.Background()

	pgC, err := containers.NewPostgres(ctx, &containers.PostgresOptions{})
	require.NoError(t, err)
	t.Cleanup(func() { errClose := pgC.Terminate(ctx); require.NoError(t, errClose) })

	getDSN, err := pgC.GetDSN(ctx)
	require.NoError(t, err)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(getDSN)))

	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	err = containers.Migrate(ctx, "../../migrations", pgC)
	require.NoError(t, err)
	return db
}

func getTables(t *testing.T, db *bun.DB) []string {
	ctx := context.Background()
	rows, err := db.QueryContext(ctx, `SELECT table_name  FROM information_schema.tables  WHERE table_schema = 'public'`)
	require.NoError(t, err)
	defer rows.Close()

	tables := make([]string, 0)
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		require.NoError(t, err)

		tables = append(tables, tableName)
	}
	return tables
}
