package logging

import (
	"github.com/TheZeroSlave/zapsentry"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func mapLoggerToSentry(logger *zap.Logger, DSN string) *zap.Logger {
	cfg := zapsentry.Configuration{
		Level: zapcore.WarnLevel, // when to send message to sentry
		Tags: map[string]string{
			"component": "system",
		},
	}
	core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromDSN(DSN))
	// in case of err it will return noop core. so we can safely attach it
	if err != nil {
		logger.Sugar().Errorf("failed to init zap", zap.Error(err))
	}

	return zapsentry.AttachCoreToLogger(core, logger)
}
