package models

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	ID          int       `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
}

type TodoCreate struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TodoUpdate struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}
