package logging

import (
	"github.com/TheZeroSlave/zapsentry"
	"go.uber.org/zap"
)

func mapLoggerToSentry(log *zap.Logger, sentryDSN string) *zap.Logger {
	cfg := zapsentry.Configuration{
		Level:             zap.WarnLevel,
		EnableBreadcrumbs: true,
		BreadcrumbLevel:   zap.WarnLevel,
		Tags: map[string]string{
			"component": "system",
		},
	}

	core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromDSN(sentryDSN))
	if err != nil {
		log.Fatal("failed to init zapsentry", zap.Error(err))
	}

	// to use breadcrumbs feature - create new scope explicitly
	log = log.With(zapsentry.NewScope())

	return zapsentry.AttachCoreToLogger(core, log)
}
