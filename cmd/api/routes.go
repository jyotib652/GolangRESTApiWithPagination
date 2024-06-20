package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) route() http.Handler {
	mux := chi.NewRouter()

	// specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))
	mux.Use(middleware.Logger)
	mux.Use(app.Authenticate)

	mux.Post("/create-employee", app.CreateEmployee)
	mux.Get("/get-employee/{id}", app.GetEmployeeByID)
	mux.Put("/update-employee/{id}", app.UpdateEmployee)
	mux.Delete("/delete-employee/{id}", app.DeleteEmployee)
	// mux.Get("/get-all-employee?{limit}=limitNumber&{cursor}=base64_string_from_previous_result", app.GetAllEmployee)
	mux.Get("/get-all-employee/{limit}/{cursor}", app.GetAllEmployee)

	return mux
}
