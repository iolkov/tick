package server

import (
	"html/template"
	"log"
	"net/http"

	"tick/internal/config"
	"tick/internal/database"
)

type application struct {
	templateCache map[string]*template.Template
	cfg           config.Config
}

type Server struct {
	*http.Server
	tmpl *template.Template
}

func New() *Server {
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

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
		cfg:           conf,
	}

	// Создаем сервер
	srv := &Server{
		Server: &http.Server{
			Handler:           app.newRouters(),
			Addr:              conf.GetServerAddress(),
			ReadTimeout:       conf.GetServerReadTimeout(),
			WriteTimeout:      conf.GetServerWriteTimeout(),
			IdleTimeout:       conf.GetServerIdleTimeout(),
			MaxHeaderBytes:    1 << 20,
			ReadHeaderTimeout: conf.GetServerReadHeaderTimeout(),
		},
	}

	srv.Handler = app.newRouters()

	return srv
}
