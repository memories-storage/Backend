package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"Backend/internal/db"
	"Backend/internal/routes"

	"github.com/joho/godotenv"
	"github.com/go-chi/chi/v5"
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

	// Basic health route
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Server is running and connected to Supabase!")
	})

	// Register auth routes
	routes.RegisterAuthRoutes(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
