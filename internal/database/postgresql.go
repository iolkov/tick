package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"tick/internal/config"

	_ "github.com/lib/pq"
)

var PDB *sql.DB

// InitDB инициализирует подключение к PostgreSQL
func InitDBPostgresql(c config.Config, l *slog.Logger) error {
	// Строка подключения
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.GetPostgresqlHost(),
		c.GetPostgresqlPort(),
		c.GetPostgresqlUser(),
		c.GetPosrgresqlPassword(),
		c.GetPostgresqlDbName(),
		c.GetPostgreSqlMode(),
	)

	var err error
	PDB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	if err = PDB.Ping(); err != nil {
		return err
	}

	if err := createTables(); err != nil {
		return err
	}

	l.Info("Database initialized",
		"status", "successfully",
		"type", "postgresql",
	)
	return nil
}

func createTables() error {
	// Таблица users (используем UUID из PostgreSQL)
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		first_name TEXT NOT NULL,
		second_name TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	// Таблица auth_providers
	authProvider := `
	CREATE TABLE IF NOT EXISTS auth_providers (
		uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		provider TEXT NOT NULL,
		name TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	// Таблица user_profiles
	userProfileProvider := `
	CREATE TABLE IF NOT EXISTS user_profiles (
		uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_uuid UUID NOT NULL,
		provider_uuid UUID NOT NULL,
		provider_user_id TEXT NOT NULL,
		full_name TEXT,
		avatar_url TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE,
		FOREIGN KEY (provider_uuid) REFERENCES auth_providers(uuid) ON DELETE CASCADE,
		UNIQUE(provider_uuid, provider_user_id)
	)`

	// Таблица todos (исправлена с INTEGER на UUID для user_id)
	todosTable := `
	CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		user_uuid UUID NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		completed BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE
	)`

	// Триггер для автоматического обновления updated_at
	triggerFunction := `
	CREATE OR REPLACE FUNCTION update_updated_at_column()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = CURRENT_TIMESTAMP;
		RETURN NEW;
	END;
	$$ language 'plpgsql'`

	// Применяем триггеры для каждой таблицы
	triggers := []string{
		`DROP TRIGGER IF EXISTS update_users_updated_at ON users`,
		`CREATE TRIGGER update_users_updated_at 
			BEFORE UPDATE ON users 
			FOR EACH ROW 
			EXECUTE FUNCTION update_updated_at_column()`,

		`DROP TRIGGER IF EXISTS update_auth_providers_updated_at ON auth_providers`,
		`CREATE TRIGGER update_auth_providers_updated_at 
			BEFORE UPDATE ON auth_providers 
			FOR EACH ROW 
			EXECUTE FUNCTION update_updated_at_column()`,

		`DROP TRIGGER IF EXISTS update_user_profiles_updated_at ON user_profiles`,
		`CREATE TRIGGER update_user_profiles_updated_at 
			BEFORE UPDATE ON user_profiles 
			FOR EACH ROW 
			EXECUTE FUNCTION update_updated_at_column()`,

		`DROP TRIGGER IF EXISTS update_todos_updated_at ON todos`,
		`CREATE TRIGGER update_todos_updated_at 
			BEFORE UPDATE ON todos 
			FOR EACH ROW 
			EXECUTE FUNCTION update_updated_at_column()`,
	}

	// Выполняем все запросы
	queries := append([]string{
		userTable,
		authProvider,
		userProfileProvider,
		todosTable,
		triggerFunction,
	}, triggers...)

	for _, query := range queries {
		if _, err := PDB.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

// func GetOrCreateUserPostgresql(firstName, secondName string) (string, error) {
// 	var userUUID string

// 	// Пытаемся найти по имени
// 	query := `SELECT uuid FROM users WHERE first_name = $1 AND second_name = $2`
// 	err := DB.QueryRow(query, firstName, secondName).Scan(&userUUID)
// 	if err == nil {
// 		return userUUID, nil
// 	}

// 	if err != sql.ErrNoRows {
// 		return "", err
// 	}

// 	// Создаем нового пользователя
// 	insertQuery := `
// 		INSERT INTO users (first_name, second_name)
// 		VALUES ($1, $2)
// 		RETURNING uuid`

// 	err = DB.QueryRow(insertQuery, firstName, secondName).Scan(&userUUID)
// 	if err != nil {
// 		return "", err
// 	}

// 	return userUUID, nil
// }
