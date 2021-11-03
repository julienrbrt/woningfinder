package logging

import (
	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
}

// NewZapLogger creates a logger using the zap library
// If sentryDSN is set errors are as well sent to Sentry
func NewZapLogger(debug bool, sentryDSN string) *Logger {
	var logger *zap.Logger
	var err error

	if debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
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
