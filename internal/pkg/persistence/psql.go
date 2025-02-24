package persistence

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

const cApplicationName = "paloma-cdp"

type Database struct {
	*bun.DB
}

func New(ctx context.Context, c *Configuration) (*Database, error) {
	pgconn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(c.Address),
		// pgdriver.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
		pgdriver.WithInsecure(true),
		pgdriver.WithUser(c.User),
		pgdriver.WithPassword(c.Password),
		pgdriver.WithDatabase(c.Database),
		pgdriver.WithApplicationName(cApplicationName),
		pgdriver.WithTimeout(c.Timeout),
	)

	sqldb := sql.OpenDB(pgconn)
	return FromSqlDB(ctx, sqldb)
}

func FromSqlDB(ctx context.Context, sqldb *sql.DB) (*Database, error) {
	db := bun.NewDB(sqldb, pgdialect.New())
	return &Database{db}, db.Ping()
}
