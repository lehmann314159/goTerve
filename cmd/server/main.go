package main

import (
	"log"
	"os"

	"github.com/mikebway/goTerve/internal/server"
	"github.com/mikebway/goTerve/internal/store"
)

func main() {
	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	// Get database path from environment or use default
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./data/terve.db"
	}

	// Initialize database
	db, err := store.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Seed data
	if err := db.Seed(); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	// Create and start server
	srv := server.New(db, port)
	log.Printf("Starting server on http://localhost:%s", port)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
