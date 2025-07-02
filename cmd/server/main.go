package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"Backend/internal/db"
	"Backend/internal/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load env
	if err := godotenv.Load(".env"); err != nil {
		_ = godotenv.Load("../../.env")
	}

	// Connect DB
	if err := db.ConnectDB(); err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer db.CloseDB()

	// Set up router
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // or "*" to allow all
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any major browsers
	}))

	// Basic health route
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Server is running and connected to Supabase!")
	})

	// Register routes
	r.Route("/api", func(api chi.Router) {
		routes.RegisterAuthRoutes(api)
		routes.RegisterImageRoutes(api)
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
