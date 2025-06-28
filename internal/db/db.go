package db

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set in environment")
	}

	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Failed to parse DB config: %v", err)
	}

	// Force IPv4 connection
	cfg.ConnConfig.DialFunc = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, "tcp4", addr)
	}

	DB, err = pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		log.Fatalf("Failed to create connection pool: %v", err)
	}

	if err := DB.Ping(context.Background()); err != nil {
		log.Fatalf("Cannot connect to DB: %v", err)
	}

	var version string
	if err := DB.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	log.Println("Connected to PostgreSQL database")
	log.Printf("Database version: %s", version)
}
