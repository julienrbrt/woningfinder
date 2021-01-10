package logging

import "go.uber.org/zap"

// NewZapLogger creates a logger using the zap library
func NewZapLogger() *zap.Logger {
	// logger, err := zap.NewProduction()
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any

	return logger
}

// NewZapLoggerWithSentry creates a logger using the zap library that throws error to Sentry
func NewZapLoggerWithSentry(sentryDSN string) *zap.Logger {
	// logger, err := zap.NewProduction()
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any

	return mapLoggerToSentry(logger, sentryDSN)
}
