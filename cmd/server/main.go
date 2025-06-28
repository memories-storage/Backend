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
	if err := godotenv.Load(".env"); err != nil {
		_ = godotenv.Load("../../.env")
	}

	// Connect to PostgreSQL (Supabase)
	if err := db.ConnectDB(); err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	// Load PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Basic test route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Server is running and connected to Supabase!")
	})

	log.Printf("Server started on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
