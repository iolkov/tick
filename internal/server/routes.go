package server

import (
	"net/http"

	"tick/internal/middleware"
	"tick/internal/server/handlers"
)

func (app application) newRouters() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", app.homeHandler)
	mux.HandleFunc("GET /api/todos", handlers.GetTodos)
	mux.HandleFunc("POST /api/todos", handlers.CreateTodo)
	mux.HandleFunc("DELETE /api/todos/{id}", handlers.DeleteTodo)
	mux.HandleFunc("PUT /api/todos/{id}", handlers.UpdateTodo)
	mux.HandleFunc("GET /health", handlers.HealthCheck(app.info))

	fileServer := http.FileServer(http.Dir("./web/"))

	mux.Handle("GET /static", http.NotFoundHandler())
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	handler := middleware.Chain(
		mux,
		middleware.RecoveryMiddleware(app.log),
		middleware.LoggingMiddleware(app.log),
		middleware.CheckDomainMiddleware(app.conf.GetDomain()),
		middleware.AuthMiddleware(app.log),
	)

	return handler
}
