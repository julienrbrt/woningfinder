package logging

import (
	"github.com/TheZeroSlave/zapsentry"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func mapLoggerToSentry(log *zap.Logger, sentryDSN string) *zap.Logger {
	cfg := zapsentry.Configuration{
		Level: zapcore.WarnLevel,
		Tags: map[string]string{
			"component": "system",
		},
	}

	core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromDSN(sentryDSN))
	if err != nil {
		log.Fatal("failed to init zapsentry", zap.Error(err))
	}

	return zapsentry.AttachCoreToLogger(core, log)
}
