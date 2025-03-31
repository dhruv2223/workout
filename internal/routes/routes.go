package routes

import (
	"github.com/dhruv2223/workout/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/health", app.HealthCheck)
	return r
}
