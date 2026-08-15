package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"tick/internal/database"
	"tick/internal/models"

	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwt := r.Header.Get("X-Auth-Request-Token")
		if jwt == "" {
			http.Error(w, "Unauthorized - No JWT token found", http.StatusUnauthorized)
			return
		}

		userClaims, err := extractClaimsFromJWT(jwt)
		if err != nil {
			fmt.Printf("Ошибка при извлечении claims: %v\n", err)
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
			fmt.Println(err)
			http.Error(w, fmt.Sprintf("Failed to get/create user: %v", err), http.StatusInternalServerError)
			return
		}

		userData.UUID = userUUID
		ctx := context.WithValue(r.Context(), "user", userData)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractClaimsFromJWT(token string) (*models.UserClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("неверный формат JWT токена")
	}

	payload, err := jwt.DecodeSegment(parts[1])
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования payload: %w", err)
	}

	var claims models.UserClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("ошибка парсинга claims: %w", err)
	}

	return &claims, nil
}
