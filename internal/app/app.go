package app

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
)

func New(ctn *Container) *App {
	return &App{
		ctn: ctn,
	}
}

type App struct {
	ctn *Container
}

func (app *App) Start(ctx context.Context) {
	srv := app.ctn.httpServer
	srv.BaseContext = func(_ net.Listener) context.Context {
		return ctx
	}

	go func() {
		app.ctn.logger.Info("starting server", zap.String("address", app.ctn.config.API().Address()))

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.ctn.logger.Fatal("start server", zap.Error(err))
		}
	}()
}

func (app *App) Shutdown(ctx context.Context) error {
	app.ctn.logger.Debug("shutting down server")

	return app.ctn.httpServer.Shutdown(ctx)
}

func (app *App) Close() error {
	return app.ctn.Close()
}
