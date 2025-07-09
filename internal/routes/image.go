package routes

import (
	"Backend/internal/handlers"
	"Backend/internal/middleware"
	
	"github.com/go-chi/chi/v5"
)

func RegisterImageRoutes(r chi.Router) {
	// Public route for adding an image
	r.Post("/upload/files", handlers.AddImageHandler)
	r.Delete("/deleteImages/{id}", handlers.DeleteImageHandler)

	r.Group(func(protected chi.Router) {
		protected.Use(middleware.AuthMiddleware)
		protected.Get("/images", handlers.GetImagesHandler)
	})
}
