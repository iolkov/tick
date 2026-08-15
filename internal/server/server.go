package server

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"tick/internal/database"
)

type application struct {
	templateCache map[string]*template.Template
}

type Server struct {
	*http.Server
	tmpl *template.Template
}

func New() *Server {
	// Инициализация БД
	if err := database.InitDB(); err != nil {
		log.Fatal(err)
	}

	templateCache, err := newTemplateCache("./web/")
	if err != nil {
		log.Fatal("failed to load templates")
	}

	app := &application{
		templateCache: templateCache,
	}

	// Создаем сервер
	srv := &Server{
		Server: &http.Server{
			Handler:           app.newRouters(),
			Addr:              getPort(),
			ReadTimeout:       60 * time.Second,
			WriteTimeout:      60 * time.Second,
			IdleTimeout:       180 * time.Second,
			MaxHeaderBytes:    1 << 20,
			ReadHeaderTimeout: 30 * time.Second,
		},
	}

	srv.Handler = app.newRouters()

	return srv
}

func getStaticDir() string {
	if dir := os.Getenv("STATIC_DIR"); dir != "" {
		return dir
	}
	return "./web"
}

func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	return ":8080"
}
