package app

import (
	"database/sql"
	"fmt"

	"github.com/kukymbr/withoutmedianews/internal/config"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	sqlDriver = "pgx"
)

func initContainer(config config.Config, logger *zap.Logger) *container {
	ctn := &container{
		config:    config,
		logger:    logger,
		finalizer: &depsFinalizer{logger: logger},
	}

	must(initDatabase(ctn), logger)

	return ctn
}

func initDatabase(ctn *container) error {
	var err error

	ctn.logger.Debug("connecting to database", zap.String("dsn", ctn.config.Db().ToDSNDebug()))

	db, err := sql.Open(sqlDriver, ctn.config.Db().ToDSN())
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	// TODO: add a retrier or something
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	ctn.db = db

	return nil
}

type container struct {
	config config.Config

	logger *zap.Logger
	db     *sql.DB

	finalizer *depsFinalizer
}

type depsFinalizer struct {
	logger *zap.Logger
	deps   []depFinalizer
}

func (d *depsFinalizer) register(name string, fn func() error) {
	d.deps = append(d.deps, depFinalizer{name, fn})
}

func (d *depsFinalizer) finalize() {
	d.logger.Debug("closing the container")

	for _, dep := range d.deps {
		d.logger.Debug("closing the dependency", zap.String("name", dep.name))

		if err := dep.fn(); err != nil {
			d.logger.Error("close dependency", zap.String("name", dep.name), zap.Error(err))
		}
	}
}

type depFinalizer struct {
	name string
	fn   func() error
}

func must(err error, logger *zap.Logger) {
	if err == nil {
		return
	}

	logger.Fatal("initialization fail", zap.Error(err))
}
