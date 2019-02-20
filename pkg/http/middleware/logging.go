package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Logging middleware
type Logging struct {
	logger *zap.Logger
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseWriter) Write(body []byte) (int, error) {
	r.size = len(body)
	return r.ResponseWriter.Write(body)
}

// NewLogging returns new Logging middleware
func NewLogging(logger *zap.Logger) *Logging {
	return &Logging{logger: logger}
}

// Middleware ...
func (l *Logging) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		l.logger.Info("",
			zap.String("method", r.Method),
			zap.String("uri", r.URL.String()),
			zap.Int("status", rw.status),
			zap.Int("size", rw.size),
			zap.Duration("duration", duration),
		)
	})
}
