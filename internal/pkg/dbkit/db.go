package dbkit

import (
	"database/sql"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	sqlDriver   = "pgx"
	goquDialect = "postgres"
)

// NewDatabase creates a new Database wrapper.
func NewDatabase(dsn string) (*Database, error) {
	db, err := sql.Open(sqlDriver, dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	return &Database{
		db:  db,
		qdb: goqu.New(goquDialect, db),
	}, nil
}

type Database struct {
	db  *sql.DB
	qdb *goqu.Database
}

func (d *Database) DB() *sql.DB {
	return d.db
}

func (d *Database) Goqu() *goqu.Database {
	return d.qdb
}

func (d *Database) Ping() error {
	if err := d.db.Ping(); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	return nil
}

func (d *Database) Close() error {
	if d == nil || d.db == nil {
		return nil
	}

	if err := d.db.Close(); err != nil {
		return fmt.Errorf("close database: %w", err)
	}

	d.qdb = nil

	return nil
}
