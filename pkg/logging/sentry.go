package logging

import (
	"github.com/TheZeroSlave/zapsentry"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func mapLoggerToSentry(logger *zap.Logger, DSN string) *zap.Logger {
	cfg := zapsentry.Configuration{
		Level:             zapcore.WarnLevel,
		EnableBreadcrumbs: true,
		BreadcrumbLevel:   zapcore.WarnLevel,
		Tags: map[string]string{
			"component": "system",
		},
	}

	core, err := zapsentry.NewCore(cfg, func() (*sentry.Client, error) {
		return sentry.NewClient(sentry.ClientOptions{
			Dsn:   DSN,
			Debug: true,
		})
	})

	// needed for breadcrumbs feature
	logger = logger.With(zapsentry.NewScope())

	// in case of err it will return noop core. so we can safely attach it
	if err != nil {
		logger.Error("failed to init zap", zap.Error(err))
	}

	return zapsentry.AttachCoreToLogger(core, logger)
}
