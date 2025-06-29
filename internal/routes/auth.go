package routes

import (
	"Backend/internal/handlers"
	"Backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterAuthRoutes(r chi.Router) {
	r.Route("/api", func(r chi.Router) {
		// Public routes
		r.Post("/signup", handlers.SignUpHandler)
		r.Post("/login", handlers.LoginHandler)
		r.Post("/forgotPassword", handlers.ForgotPasswordHandler)
		r.Post("/resetPassword", handlers.ResetPasswordHandler)

		// Protected routes
		r.Group(func(protected chi.Router) {
			protected.Use(middleware.AuthMiddleware)

			// Authenticated user data
			protected.Get("/getUserData", handlers.GetUserHandler)
			protected.Post("/changePassword", handlers.ChangePasswordHandler)
			protected.Post("/changeUserProfile",handlers.UpdateUserProfileHandler)
		})
	})
}
