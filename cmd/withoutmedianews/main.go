package main

import (
	"context"
	"log"

	"github.com/kukymbr/withoutmedianews/internal/app"
	"github.com/kukymbr/withoutmedianews/internal/config"
	"github.com/kukymbr/withoutmedianews/internal/pkg/logkit"
	"go.uber.org/zap"
)

func main() {
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	logger, closeLogger := logkit.New(conf.Logger().Level())
	defer closeLogger()

	logger.Debug("config loaded", zap.Any("conf", conf.DebugJSON()))

	app.Run(context.Background(), conf, logger)
}
