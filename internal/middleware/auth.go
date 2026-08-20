package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"tick/internal/database"
	"tick/internal/models"

	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("X-Auth-Request-Token")
			if tokenString == "" {
				http.Error(w, "Unauthorized - No JWT token found", http.StatusUnauthorized)
				return
			}

			userClaims, err := extractClaimsFromJWT(tokenString)
			if err != nil {
				log.Error("Error extracting JWT claims", "error", err)
				http.Error(w, "Unauthorized - Invalid token", http.StatusUnauthorized)
				return
			}

			userData := &models.UserData{
				DisplayName: userClaims.Name,
				Email:       userClaims.Email,
				ProviderID:  userClaims.Sub,
				ISS:         userClaims.ISS,
				AvatarURL:   userClaims.Picture,
			}

			userUUID, err := database.GetUser(userData)
			if err != nil {
				log.Error("Failed to get/create user", "error", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			userData.UUID = userUUID
			ctx := context.WithValue(r.Context(), "user", userData)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractClaimsFromJWT(token string) (*models.UserClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT token format")
	}

	payload, err := jwt.DecodeSegment(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %w", err)
	}

	var claims models.UserClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	return &claims, nil
}
