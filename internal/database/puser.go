package database

import (
	"database/sql"
	"fmt"
	"strings"
	"tick/internal/models"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func PGetUser(user *models.UserData) (uuid.UUID, error) {
	tx, err := PDB.Begin()
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
		WHERE ap.provider = $1 AND up.provider_user_id = $2
	`

	err = PDB.QueryRow(query, user.ISS, user.ProviderID).Scan(&userUUID)
	if err == nil {
		return userUUID, nil
	}

	if err != sql.ErrNoRows {
		return uuid.Nil, err
	}

	userUUID, err = pcreateUser(user)
	if err != nil {
		return uuid.Nil, err
	}
	return userUUID, nil
}

func pcreateUser(user *models.UserData) (uuid.UUID, error) {
	tx, err := PDB.Begin()
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback()

	providerUUID, err := pgetOrCreateProvider(tx, user.ISS)
	if err != nil {
		return uuid.Nil, fmt.Errorf("ошибка с провайдером: %w", err)
	}

	userUUID := uuid.New()
	parts := strings.Split(user.DisplayName, " ")
	if len(parts) != 2 {
		return uuid.Nil, fmt.Errorf("имя должно состоять из двух слов: %s", user.DisplayName)
	}

	_, err = tx.Exec(`
        INSERT INTO users (uuid, first_name, second_name) 
        VALUES ($1, $2, $3)
    `, userUUID.String(), parts[0], parts[1])
	if err != nil {
		return uuid.Nil, fmt.Errorf("ошибка создания пользователя: %w", err)
	}

	profileUUID := uuid.New()
	_, err = tx.Exec(`
        INSERT INTO user_profiles (
            uuid, user_uuid, provider_uuid, provider_user_id, 
            full_name, avatar_url
        ) VALUES ($1, $2, $3, $4, $5, $6)
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

func pgetOrCreateProvider(tx *sql.Tx, providerName string) (uuid.UUID, error) {
	var providerUUIDStr string
	err := tx.QueryRow(`SELECT uuid FROM auth_providers WHERE provider = $1`, providerName).Scan(&providerUUIDStr)

	if err == nil {
		return uuid.Parse(providerUUIDStr)
	}

	if err != sql.ErrNoRows {
		return uuid.Nil, err
	}

	newUUID := uuid.New().String()
	_, err = tx.Exec(`
        INSERT INTO auth_providers (uuid, provider, name, created_at, updated_at)
        VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
    `, newUUID, providerName, providerName)
	if err != nil {
		return uuid.Nil, fmt.Errorf("ошибка вставки провайдера: %w", err)
	}

	return uuid.Parse(newUUID)
}
