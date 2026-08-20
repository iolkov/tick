package handlers

import (
	"encoding/json"
	"net/http"
	"tick/internal/info"
	"time"
)

func HealthCheck(info info.AppInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"status":  "ok",
			"service": "tick",
			"time":    time.Now().UTC(),
			"version": info.Version,
			"build":   info.Build,
			"date":    info.Date,
			"start":   info.Start,
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")

		json.NewEncoder(w).Encode(response)
	}
}
