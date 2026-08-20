package database

import (
	"database/sql"
	"log/slog"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(file string, log *slog.Logger) error {
	var err error
	DB, err = sql.Open("sqlite", file)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	userTable := `
    CREATE TABLE IF NOT EXISTS users (
        uuid TEXT PRIMARY KEY,
        first_name TEXT NOT NULL,
        second_name TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )`

	authProvider := `
    CREATE TABLE IF NOT EXISTS auth_providers (
        uuid TEXT PRIMARY KEY,
        provider TEXT NOT NULL,
        name TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )`

	userProfileProvider := `
    CREATE TABLE IF NOT EXISTS user_profiles (
        uuid TEXT PRIMARY KEY,
        user_uuid TEXT NOT NULL,
        provider_uuid TEXT NOT NULL,
        provider_user_id TEXT NOT NULL,
        full_name TEXT,
        avatar_url TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE,
        FOREIGN KEY (provider_uuid) REFERENCES auth_providers(uuid) ON DELETE CASCADE,
        UNIQUE(provider_uuid, provider_user_id)
    )`

	todosTable := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		completed BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
	)`

	indexQuery := `CREATE INDEX IF NOT EXISTS idx_todos_user_id ON todos(user_id)`

	if _, err := DB.Exec(todosTable); err != nil {
		return err
	}
	if _, err := DB.Exec(indexQuery); err != nil {
		return err
	}
	if _, err := DB.Exec(userTable); err != nil {
		return err
	}
	if _, err := DB.Exec(userProfileProvider); err != nil {
		return err
	}
	if _, err := DB.Exec(authProvider); err != nil {
		return err
	}

	log.Info("Database initialized",
		"status", "successfully",
	)
	return nil
}
func GetOrCreateUser(username, email string) (int64, error) {
	var userID int64

	err := DB.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err == nil {
		return userID, nil
	}

	if email != "" {
		err = DB.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&userID)
		if err == nil {
			return userID, nil
		}
	}

	result, err := DB.Exec("INSERT INTO users (username, email) VALUES (?, ?)", username, email)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return id, nil
}
