package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() error {
	// Load DB URL from environment variable
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL not set in environment")
	}

	// Open the database connection
	var err error
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("failed to open DB connection: %v", err)
	}

	// Verify the connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("cannot connect to DB: %v", err)
	}

	fmt.Println("Connected to PostgreSQL database")
	return nil
}
