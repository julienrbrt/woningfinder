package logging

import (
	"context"

	"go.uber.org/zap"
)

// Logger embbed zap.Logger
type Logger struct {
	*zap.Logger
}

// Printf prints a log as info
func (l *Logger) Printf(_ context.Context, template string, args ...interface{}) {
	l.Sugar().Infof(template, args)
}

// NewZapLogger creates a logger using the zap library
// If debug is false errors are as well sent to Sentry
func NewZapLogger(debug bool, sentryDSN string) *Logger {
	var err error
	var logger *zap.Logger

	// debug logger
	if debug {
		logger, err = zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
		defer logger.Sync() // flushes buffer, if any

		return &Logger{logger}
	}

	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	return &Logger{mapLoggerToSentry(logger, sentryDSN)}
}

// NewZapLoggerWithoutSentry default the NewZapLogger without Sentry
func NewZapLoggerWithoutSentry() *Logger {
	return NewZapLogger(true, "")
}
