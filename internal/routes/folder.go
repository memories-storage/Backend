package routes

import (
	"Backend/internal/handlers"
	"Backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterFolderRoutes(r chi.Router) {
	r.Group(func(protected chi.Router) {
		protected.Use(middleware.AuthMiddleware)
		
		// Folder management endpoints
		protected.Post("/folders", handlers.CreateFolderHandler)
		protected.Get("/folders", handlers.GetFoldersHandler)
		protected.Delete("/folders/{id}", handlers.DeleteFolderHandler)
	})
} 