package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"Backend/internal/db"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load("../../.env"); err != nil {
		log.Printf("No .env file found or failed to load: %v", err)
	}

	// Connect to DB
	db.ConnectDB()

	// Load port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Basic route handler (you'll wire routes later)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Server is running!")
	})

	// Start the server
	log.Printf("Server started on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
