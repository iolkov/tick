package database

import (
	"database/sql"
	"fmt"
	"strings"
	"tick/internal/models"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func GetUser(user *models.UserData) (uuid.UUID, error) {
	tx, err := DB.Begin()
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback()

	var userUUID uuid.UUID

	query := `
		SELECT u.uuid 
		FROM users u
		JOIN user_profiles up ON u.uuid = up.user_uuid
		JOIN auth_providers ap ON up.provider_uuid = ap.uuid
		WHERE ap.provider = ? AND up.provider_user_id = ?
	`

	err = DB.QueryRow(query, user.ISS, user.ProviderID).Scan(&userUUID)
	if err == nil {
		return userUUID, nil
	}

	if err != sql.ErrNoRows {
		return uuid.Nil, err
	}

	userUUID, err = createUser(user)
	if err != nil {
		return uuid.Nil, err
	}
	return userUUID, nil
}

func createUser(user *models.UserData) (uuid.UUID, error) {
	tx, err := DB.Begin()
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback()

	// 1. Получаем или создаем provider
	providerUUID, err := getOrCreateProvider(tx, user.ISS)
	if err != nil {
		return uuid.Nil, fmt.Errorf("ошибка с провайдером: %w", err)
	}

	// 2. Создаем пользователя
	userUUID := uuid.New()
	parts := strings.Split(user.DisplayName, " ")
	if len(parts) != 2 {
		return uuid.Nil, fmt.Errorf("имя должно состоять из двух слов: %s", user.DisplayName)
	}

	_, err = tx.Exec(`
        INSERT INTO users (uuid, first_name, second_name) 
        VALUES (?, ?, ?)
    `, userUUID.String(), parts[0], parts[1])
	if err != nil {
		return uuid.Nil, fmt.Errorf("ошибка создания пользователя: %w", err)
	}

	// 3. Создаем профиль
	profileUUID := uuid.New()
	_, err = tx.Exec(`
        INSERT INTO user_profiles (
            uuid, user_uuid, provider_uuid, provider_user_id, 
            full_name, avatar_url
        ) VALUES (?, ?, ?, ?, ?, ?)
    `, profileUUID.String(), userUUID.String(), providerUUID.String(),
		user.ProviderID, user.DisplayName, user.AvatarURL)
	if err != nil {
		return uuid.Nil, fmt.Errorf("ошибка создания профиля: %w", err)
	}

	user.UUID = userUUID

	if err = tx.Commit(); err != nil {
		return uuid.Nil, err
	}

	return userUUID, nil
}

// Единая функция для получения или создания провайдера
func getOrCreateProvider(tx *sql.Tx, providerName string) (uuid.UUID, error) {
	var providerUUIDStr string
	err := tx.QueryRow(`SELECT uuid FROM auth_providers WHERE provider = ?`, providerName).Scan(&providerUUIDStr)

	if err == nil {
		// Нашли существующий
		return uuid.Parse(providerUUIDStr)
	}

	if err != sql.ErrNoRows {
		return uuid.Nil, err
	}

	// Создаем новый
	newUUID := uuid.New().String()
	_, err = tx.Exec(`
        INSERT INTO auth_providers (uuid, provider, name, created_at, updated_at)
        VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
    `, newUUID, providerName, providerName)
	if err != nil {
		return uuid.Nil, fmt.Errorf("ошибка вставки провайдера: %w", err)
	}

	return uuid.Parse(newUUID)
}
