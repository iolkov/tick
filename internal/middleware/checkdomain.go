package middleware

import (
	"net/http"
)

func CheckDomainMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Forwarded-Host") != "tick.iolkov.ru" {
			w.WriteHeader(423)
			return
		}

		next.ServeHTTP(w, r)
	})
}
