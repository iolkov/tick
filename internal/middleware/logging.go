package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &responseWriterWrapper{
				ResponseWriter: w,
				statusCode:     0,
			}

			next.ServeHTTP(wrapped, r)

			// if wrapped.statusCode == 0 {
			// 	wrapped.statusCode = http.StatusOK
			// }

			log.Info("request",
				"ip", r.RemoteAddr,
				"x-forwarded-proto", r.Header.Get("X-Forwarded-Proto"),
				"x-forwarded-host", r.Header.Get("X-Forwarded-Host"),
				"x-forwarded-port", r.Header.Get("X-Forwarded-Port"),
				"x-forwarded-for", r.Header.Get("X-Forwarded-For"),
				"x-request-id", r.Header.Get("X-Request-Id"),
				"x-real-ip", r.Header.Get("X-Real-Ip"),
				"method", r.Method,
				"path", r.URL.Path,
				"status", wrapped.statusCode,
				"duration", time.Since(start).Milliseconds(),
			)
		})
	}
}
