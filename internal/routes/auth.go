package routes

import (
	"Backend/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func RegisterAuthRoutes(r chi.Router) {
	r.Route("/api", func(r chi.Router) {
		r.Post("/signup", handlers.SignUpHandler)
		r.Post("/login", handlers.LoginHandler)
	})
}
