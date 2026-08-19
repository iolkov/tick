package server

import (
	"html/template"
	"log/slog"
	"net/http"

	"tick/internal/config"
	"tick/internal/database"
	"tick/internal/info"
)

type application struct {
	templateCache map[string]*template.Template
	cfg           config.Config
	info          info.AppInfo
	log           *slog.Logger
}

type Server struct {
	*http.Server
	tmpl *template.Template
}

func New(info *info.AppInfo, conf config.Config, log *slog.Logger) *Server {
	if err := database.InitDB(conf.GetDatabase(), log); err != nil {
		log.Error("initialized database",
			"error", err,
		)
	}

	templateCache, err := newTemplateCache(conf.GetTemplateDir())
	if err != nil {
		log.Error("failed to load templates",
			"error", err,
		)
	}

	app := &application{
		templateCache: templateCache,
		cfg:           conf,
		info:          *info,
		log:           log,
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

	return srv
}
