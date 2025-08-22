package middleware

import (
	"net/http"
	"time"

	"github.com/4Noyis/my-library/internal/logger"
	"github.com/sirupsen/logrus"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += int64(n)
	return n, err
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the response writer to capture status code
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     200, // default status code
		}

		// Log the incoming request
		logger.Logger.WithFields(logrus.Fields{
			"method":      r.Method,
			"path":        r.URL.Path,
			"remote_addr": r.RemoteAddr,
			"user_agent":  r.UserAgent(),
			"type":        "request_start",
		}).Info("Request started")

		// Call the next handler
		next.ServeHTTP(wrapped, r)

		// Calculate duration
		duration := time.Since(start)

		// Log the completed request
		logger.Logger.WithFields(logrus.Fields{
			"method":        r.Method,
			"path":          r.URL.Path,
			"remote_addr":   r.RemoteAddr,
			"status_code":   wrapped.statusCode,
			"duration_ms":   duration.Milliseconds(),
			"bytes_written": wrapped.written,
			"type":          "request_complete",
		}).Info("Request completed")
	})
}
