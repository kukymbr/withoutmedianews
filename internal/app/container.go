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
	"github.com/kukymbr/withoutmedianews/internal/domain"
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
	initServices(ctn)
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
	ctn.newsRepo = db.NewNewsRepository(ctn.db.DB())
}

func initServices(ctn *Container) {
	ctn.newsService = domain.NewNewsService(ctn.newsRepo)
}

func initServer(ctn *Container, ctx context.Context) {
	ctn.errResponder = server.NewErrorResponder(ctn.logger)

	ctn.server = &server.Server{
		NewsController:       controller.NewNewsController(ctn.newsService),
		CategoriesController: controller.NewCategoriesController(ctn.newsService),
		TagsController:       controller.NewTagsController(ctn.newsService),
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

	newsRepo    *db.NewsRepository
	newsService *domain.Service

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
