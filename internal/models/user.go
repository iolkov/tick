package models

import (
	"time"

	"github.com/google/uuid"
)

type UserData struct {
	DisplayName string    `json:"name"`
	Email       string    `json:"email"`
	UUID        uuid.UUID `json:"uuis"`
	ProviderID  string    `json:"sub"`
	ISS         string    `json:"iss"`
	AvatarURL   string    `json:"avatar_url"`
}

type UserClaims struct {
	Sub      string `json:"sub"`
	ISS      string `json:"iss"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Profile  string `json:"profile"`
	Picture  string `json:"picture"`
}

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// UserProfile - расширенная информация о пользователе
type UserProfile struct {
	ID         uuid.UUID `json:"id" db:"id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	FullName   string    `json:"full_name" db:"full_name"`
	AvatarURL  string    `json:"avatar_url" db:"avatar_url"`
	ProfileURL string    `json:"profile_url" db:"profile_url"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// AuthProvider - информация об аккаунте у провайдера
type AuthProvider struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	Provider     string    `json:"provider" db:"provider"`
	ProviderID   string    `json:"provider_id" db:"provider_id"`
	AccessToken  string    `json:"access_token,omitempty" db:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty" db:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
