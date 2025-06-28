// package db

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"os"

// 	_ "github.com/lib/pq"
// )

// var DB *sql.DB

// func ConnectDB() {
// 	// Load DB URL from environment variable
// 	dbURL := os.Getenv("DATABASE_URL")
// 	if dbURL == "" {
// 		log.Fatal("DATABASE_URL not set in environment")
// 	}

// 	// Open the database connection
// 	var err error
// 	DB, err = sql.Open("postgres", dbURL)
// 	if err != nil {
// 		log.Fatalf("Failed to open DB connection: %v", err)
// 	}

// 	// Verify the connection
// 	if err := DB.Ping(); err != nil {
// 		log.Fatalf("Cannot connect to DB: %v", err)
// 	}

// 	fmt.Println("Connected to PostgreSQL database")
// }




package db

import (
	"context"
	"log"
	"os"
	"github.com/jackc/pgx/v5"
)

func ConnectDB() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	// Example query to test connection
	var version string
	if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	log.Println("Connected to:", version)
}