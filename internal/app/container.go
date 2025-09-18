package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/kukymbr/retrier"
	"github.com/kukymbr/withoutmedianews/internal/api/http/controller"
	"github.com/kukymbr/withoutmedianews/internal/api/http/server"
	"github.com/kukymbr/withoutmedianews/internal/config"
	"github.com/kukymbr/withoutmedianews/internal/db"
	"github.com/kukymbr/withoutmedianews/internal/pkg/dbkit"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func BuildContainer(
	backgroundCtx context.Context,
	requestsCtx context.Context,
	config config.Config,
	logger *zap.Logger,
) (*Container, error) {
	ctn := &Container{
		config:    config,
		logger:    logger,
		finalizer: &depsFinalizer{logger: logger},
		retrier:   retrier.NewLinear(5, 10*time.Second),
	}

	if err := initDatabase(ctn, backgroundCtx); err != nil {
		return nil, err
	}

	initRepositories(ctn)
	initServer(ctn, requestsCtx)

	return ctn, nil
}

func initDatabase(ctn *Container, ctx context.Context) error {
	var err error

	ctn.logger.Debug("connecting to database", zap.String("dsn", ctn.config.Db().ToDSNDebug()))

	database, err := dbkit.NewDatabase(ctn.config.Db())
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	if err := ctn.retrier.Do(ctx, func() error {
		return database.Ping(ctx)
	}); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	ctn.db = database

	ctn.finalizer.register("database", database.Close)

	return nil
}

func initRepositories(ctn *Container) {
	repo := db.NewWithoutmedianewsRepo(ctn.db.DB())
	ctn.newsRepo = &repo
}

func initServer(ctn *Container, ctx context.Context) {
	ctn.errResponder = server.NewErrorResponder(ctn.logger)

	ctn.server = &server.Server{
		NewsController:       controller.NewNewsController(ctn.newsRepo),
		CategoriesController: controller.NewCategoriesController(ctn.newsRepo),
		TagsController:       controller.NewTagsController(ctn.newsRepo),
	}

	ctn.router = initRouter(ctn.server, ctn.errResponder)

	ctn.httpServer = &http.Server{
		Addr:              ctn.config.API().Address(),
		Handler:           ctn.router,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}
}

type Container struct {
	config config.Config

	logger *zap.Logger
	db     *dbkit.Database

	newsRepo *db.WithoutmedianewsRepo

	errResponder *server.ErrorResponder
	router       http.Handler
	httpServer   *http.Server
	server       *server.Server

	retrier   retrier.Retrier
	finalizer *depsFinalizer
}

func (ctn *Container) GetRouter() http.Handler {
	return ctn.router
}

func (ctn *Container) Close() error {
	ctn.finalizer.finalize()

	return nil
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
