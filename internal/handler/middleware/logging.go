package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/woningfinder/woningfinder/pkg/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// CreateZapMiddleware logs requests
func CreateZapMiddleware(logger *logging.Logger) func(next http.Handler) http.Handler {
	return routerLogger{logger: logger}.middleware
}

type routerLogger struct {
	logger *logging.Logger
}

func (m routerLogger) middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := middleware.GetReqID(r.Context())
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		latency := time.Since(start)

		fields := []zapcore.Field{
			zap.String("url", r.URL.String()),
			zap.String("uri", r.RequestURI),
			zap.Duration("took", latency),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Int("http_status", ww.Status()),
			zap.String("http_method", r.Method),
			zap.String("user_agent", r.UserAgent()),
		}
		if requestID != "" {
			fields = append(fields, zap.String("request_id", requestID))
		}

		// checks the status code of the response and logs are error if server error
		logLine := fmt.Sprintf("%s %s", r.Method, r.URL.String())
		if ww.Status() >= http.StatusInternalServerError {
			m.logger.Error(fmt.Sprintf("%s - Invalid response with status code %d", logLine, ww.Status()), fields...)
		} else {
			m.logger.Info(logLine, fields...)
		}

	}

	return http.HandlerFunc(fn)
}
