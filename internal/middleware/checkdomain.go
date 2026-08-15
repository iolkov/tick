package middleware

import (
	"net/http"
)

func CheckDomainMiddleware(domain string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Если заголовок не совпадает - возвращаем 423
			if r.Header.Get("X-Forwarded-Host") != domain {
				http.Error(w, "Forbidden", http.StatusLocked)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
