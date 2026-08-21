package database

import (
	"tick/internal/models"

	"github.com/google/uuid"
)

func GetTodosByUserID(userUUID uuid.UUID) ([]models.Todo, error) {
	query := `
		SELECT id, title, description, completed, created_at
		FROM todos 
		WHERE user_uuid = $1 
		ORDER BY 
			CASE WHEN completed = false THEN 0 ELSE 1 END,
			created_at DESC
	`

	rows, err := PDB.Query(query, userUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []models.Todo{}
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
			&todo.CreatedAt,
		)
		if err != nil {
			continue
		}
		todo.UserID = userUUID
		todos = append(todos, todo)
	}

	return todos, nil
}

func PCreateTodo(todo *models.Todo) error {
	query := `
		INSERT INTO todos (user_uuid, title, description) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at
	`

	err := PDB.QueryRow(
		query,
		todo.UserID,
		todo.Title,
		todo.Description,
	).Scan(&todo.ID, &todo.CreatedAt)

	return err
}
