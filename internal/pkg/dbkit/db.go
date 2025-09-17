package dbkit

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kukymbr/withoutmedianews/internal/config"
)

// NewDatabase creates a new Database wrapper.
func NewDatabase(conf config.DbConfig) (*Database, error) {
	db := pg.Connect(&pg.Options{
		Addr:     conf.Address(),
		User:     conf.Username(),
		Password: conf.Password(),
		Database: conf.Database(),
	})

	return &Database{
		db: db,
	}, nil
}

type Database struct {
	db *pg.DB
}

func (d *Database) DB() *pg.DB {
	return d.db
}

func (d *Database) Ping(ctx context.Context) error {
	if err := d.db.Ping(ctx); err != nil {
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

	d.db = nil

	return nil
}
