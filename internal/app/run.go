package app

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/kukymbr/withoutmedianews/internal/config"
	"go.uber.org/zap"
)

const (
	shutdownTimeout   = 15 * time.Second
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
)

func Run(backgroundCtx context.Context, conf config.Config, logger *zap.Logger) {
	signalCtx, stop := signal.NotifyContext(backgroundCtx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	requestsCtx, stopRequestsGracefully := context.WithCancel(backgroundCtx)

	defer func() {
		r := recover()
		if r == nil {
			return
		}

		logger.Panic(
			"application panic occurred",
			zap.Any("error", r), zap.Stack("stack"),
		)
	}()

	ctn := initContainer(backgroundCtx, conf, logger)
	defer ctn.finalizer.finalize()

	server := &http.Server{
		Addr:              conf.API().Address(),
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		BaseContext: func(_ net.Listener) context.Context {
			return requestsCtx
		},
		Handler: ctn.router,
	}

	go func() {
		logger.Info("starting server", zap.String("address", conf.API().Address()))

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("start server", zap.Error(err))
		}
	}()

	<-signalCtx.Done()
	stop()

	shutdownCtx, cancel := context.WithTimeout(backgroundCtx, shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error", zap.Error(err))
	}

	stopRequestsGracefully()
}
