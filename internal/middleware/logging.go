package middleware

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		logger.Info("request",
			"ip", r.RemoteAddr,
			"x-forwarded-proto", r.Header.Get("X-Forwarded-Proto"),
			"x-forwarded-host", r.Header.Get("X-Forwarded-Host"),
			"x-forwarded-port", r.Header.Get("X-Forwarded-Port"),
			"x-forwarded-for", r.Header.Get("X-Forwarded-For"),
			"x-real-ip", r.Header.Get("X-Real-Ip"),
			"method", r.Method,
			"path", r.URL.Path,
			"query_params", r.URL.RawQuery,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}
