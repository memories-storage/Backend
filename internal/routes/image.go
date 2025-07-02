package routes

import (
	"Backend/internal/handlers"
	"github.com/go-chi/chi/v5"
	"Backend/internal/middleware"
)

func RegisterImageRoutes(r chi.Router) {
	// Public route for adding an image
	r.Post("/addImages/{userId}", handlers.AddImageHandler)
	r.Delete("/deleteImages/{id}", handlers.DeleteImageHandler)

	r.Group(func(protected chi.Router) {
		protected.Use(middleware.AuthMiddleware)
		protected.Get("/images", handlers.GetImagesHandler)
	})
}
