package logkit

import (
	"os"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

// New creates new console logger instance.
func New(level zapcore.Level) (logger *zap.Logger, closeFn func()) {
	// TODO: environments support.
	config := zap.NewProductionEncoderConfig()

	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncodeName = zapcore.FullNameEncoder
	config.NameKey = "who"

	consoleEncoder := zapcore.NewConsoleEncoder(config)

	stdErrLevel := func(lvl zapcore.Level) bool {
		return lvl >= level && lvl >= zapcore.ErrorLevel
	}
	stdOutLevel := func(lvl zapcore.Level) bool {
		return lvl >= level && lvl < zapcore.ErrorLevel
	}

	stdout := zapcore.Lock(os.Stdout)
	stderr := zapcore.Lock(os.Stderr)

	cores := make([]zapcore.Core, 0, 2)

	cores = append(cores, zapcore.NewCore(consoleEncoder, stdout, zap.LevelEnablerFunc(stdOutLevel)))
	cores = append(cores, zapcore.NewCore(consoleEncoder, stderr, zap.LevelEnablerFunc(stdErrLevel)))

	core := zapcore.NewTee(cores...)
	logger = zap.New(core)

	closeFn = func() {
		_ = logger.Sync()
	}

	return logger, closeFn
}

// NewTestLogger creates logger for tests
func NewTestLogger(t *testing.T) *zap.Logger {
	logger := zaptest.NewLogger(t)

	t.Cleanup(func() {
		_ = logger.Sync()
	})

	return logger
}
