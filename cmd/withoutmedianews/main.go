package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/kukymbr/withoutmedianews/internal/app"
	"github.com/kukymbr/withoutmedianews/internal/config"
	"github.com/kukymbr/withoutmedianews/internal/pkg/logkit"
	"go.uber.org/zap"
)

const shutdownTimeout = 15 * time.Second

func main() {
	backgroundCtx := context.Background()

	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	logger, closeLogger := logkit.New(conf.Logger().Level())
	defer closeLogger()

	logger.Debug("config loaded", zap.Any("conf", conf.DebugJSON()))

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

	ctn, err := app.BuildContainer(backgroundCtx, requestsCtx, conf, logger)
	if err != nil {
		logger.Fatal("failed to build container", zap.Error(err))
	}

	application := app.New(ctn)
	defer func() {
		_ = application.Close()
	}()

	application.Start(requestsCtx)

	<-signalCtx.Done()
	stop()

	shutdownCtx, cancel := context.WithTimeout(backgroundCtx, shutdownTimeout)
	defer cancel()

	if err := application.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error", zap.Error(err))
	}

	stopRequestsGracefully()
}
