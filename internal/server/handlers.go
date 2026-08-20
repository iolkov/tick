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

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *TemplateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.log.Error("Шаблон не существует", "template", name)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err := ts.Execute(w, td)
	if err != nil {
		app.log.Error("Ошибка поиска шаблона", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (app application) homeHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := handlers.GetUserFromContext(r)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	if r.URL.Path != "/" {
		http.NotFound(w, r)
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
