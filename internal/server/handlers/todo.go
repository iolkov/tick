package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"tick/internal/database"

	"tick/internal/models"
)

func GetUserFromContext(r *http.Request) (*models.UserData, bool) {
	user, ok := r.Context().Value("user").(*models.UserData)
	return user, ok
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	rows, err := database.DB.Query(
		`SELECT id, title, description, completed, created_at 
			FROM todos 
			WHERE user_id = ? 
			ORDER BY 
			CASE WHEN completed = 0 THEN 0 ELSE 1 END,
			created_at DESC`,
		user.UUID,
	)
	if err != nil {
		http.Error(w, "Failed to fetch todos: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	todos := []models.Todo{}
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description,
			&todo.Completed, &todo.CreatedAt)
		if err != nil {
			continue
		}
		todo.UserID = user.UUID
		todos = append(todos, todo)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var req models.TodoCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validation
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec(
		"INSERT INTO todos (user_id, title, description) VALUES (?, ?, ?)",
		user.UUID, req.Title, req.Description,
	)
	if err != nil {
		http.Error(w, "Failed to create todo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	// Get created todo
	var todo models.Todo
	err = database.DB.QueryRow(
		`SELECT id, title, description, completed, created_at
		 FROM todos
		 WHERE id = ?`,
		id,
	).Scan(&todo.ID, &todo.Title, &todo.Description,
		&todo.Completed, &todo.CreatedAt)

	if err != nil {
		http.Error(w, "Failed to fetch created todo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	todo.UserID = user.UUID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Extract ID from URL
	path := strings.TrimPrefix(r.URL.Path, "/api/todos/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec(
		"DELETE FROM todos WHERE id = ? AND user_id = ?",
		id, user.UUID,
	)
	if err != nil {
		http.Error(w, "Failed to delete todo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Todo deleted successfully",
		"id":      id,
	})
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Extract ID from URL
	path := strings.TrimPrefix(r.URL.Path, "/api/todos/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var req models.TodoUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if todo exists
	var exists bool
	err = database.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM todos WHERE id = ? AND user_id = ?)",
		id, user.UUID,
	).Scan(&exists)

	if err != nil || !exists {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	_, err = database.DB.Exec(
		`UPDATE todos
		 SET title = ?, description = ?, completed = ?
		 WHERE id = ? AND user_id = ?`,
		req.Title, req.Description, req.Completed, id, user.UUID,
	)
	if err != nil {
		http.Error(w, "Failed to update todo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get updated todo
	var todo models.Todo
	err = database.DB.QueryRow(
		`SELECT id, title, description, completed, created_at
		 FROM todos
		 WHERE id = ?`,
		id,
	).Scan(&todo.ID, &todo.Title, &todo.Description,
		&todo.Completed, &todo.CreatedAt)

	if err != nil {
		http.Error(w, "Failed to fetch updated todo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	todo.UserID = user.UUID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}
