package server

import (
	"net/http"
	"tick/internal/server/handlers"

	"github.com/google/uuid"
)

type TemplateData struct {
	Email     string
	UserID    uuid.UUID
	Name      string
	AvatarURL string
}

func (app application) homeHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := handlers.GetUserFromContext(r)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	if r.URL.Path != "/" {
		http.NotFound(w, r) // Отправляет стандартную 404 страницу
		return
	}

	data := TemplateData{
		Email:     user.Email,
		Name:      user.DisplayName,
		UserID:    user.UUID,
		AvatarURL: user.AvatarURL,
	}

	app.render(w, r, "home.page.html", &data)
}

func (app application) templateHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := handlers.GetUserFromContext(r)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	data := TemplateData{
		Email:     user.Email,
		Name:      user.DisplayName,
		UserID:    user.UUID,
		AvatarURL: user.AvatarURL,
	}

	app.render(w, r, "template.page.html", &data)
}
