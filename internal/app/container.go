package app

import (
	"fmt"

	apihttp "github.com/kukymbr/withoutmedianews/internal/api/http"
	"github.com/kukymbr/withoutmedianews/internal/config"
	"github.com/kukymbr/withoutmedianews/internal/news"
	"github.com/kukymbr/withoutmedianews/internal/news/repository"
	"github.com/kukymbr/withoutmedianews/internal/pkg/dbkit"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func initContainer(config config.Config, logger *zap.Logger) *container {
	ctn := &container{
		config:    config,
		logger:    logger,
		finalizer: &depsFinalizer{logger: logger},
	}

	must(initDatabase(ctn), logger)

	initRepositories(ctn)
	initServices(ctn)
	initServer(ctn)

	return ctn
}

func initDatabase(ctn *container) error {
	var err error

	ctn.logger.Debug("connecting to database", zap.String("dsn", ctn.config.Db().ToDSNDebug()))

	db, err := dbkit.NewDatabase(ctn.config.Db().ToDSN())
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	// TODO: add a retrier or something
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	ctn.db = db

	ctn.finalizer.register("database", db.Close)

	return nil
}

func initRepositories(ctn *container) {
	ctn.newsRepo = repository.NewNewsRepository(ctn.db, ctn.logger)
}

func initServices(ctn *container) {
	ctn.newsService = news.NewNewsService(ctn.newsRepo)
}

func initServer(ctn *container) {
	ctn.server = &apihttp.Server{
		NewsController: apihttp.NewNewsController(ctn.newsService),
	}
}

type container struct {
	config config.Config

	logger *zap.Logger
	db     *dbkit.Database

	newsRepo    *repository.NewsRepository
	newsService *news.News
	server      *apihttp.Server

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
