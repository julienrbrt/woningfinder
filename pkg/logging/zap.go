package logging

import "go.uber.org/zap"

// NewZapLogger creates a logger using the zap library
// If debug is true errors are as well sent to Sentry
func NewZapLogger(debug bool, sentryDSN string) *zap.Logger {
	// logger, err := zap.NewProduction()
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any

	if debug {
		return logger
	}

	return mapLoggerToSentry(logger, sentryDSN)
}

// NewTestZapLogger default the NewZapLogger without Sentry
func NewTestZapLogger() *zap.Logger {
	return NewZapLogger(false, "")
}
