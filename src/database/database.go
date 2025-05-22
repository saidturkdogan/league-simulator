package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/user/league-simulator/src/config"
)

// Connect establishes a connection to the database
func Connect(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.GetDatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Check database connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// RunMigrations runs database migrations
func RunMigrations(db *sql.DB) error {
	log.Println("Running database migrations...")

	// Read schema SQL file
	schemaSQL, err := os.ReadFile("src/database/migrations/schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Execute schema SQL
	_, err = db.Exec(string(schemaSQL))
	if err != nil {
		return fmt.Errorf("failed to execute schema migration: %w", err)
	}

	// Verify migrations by checking if tables exist
	var tableCount int
	err = db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tableCount)
	if err != nil {
		return fmt.Errorf("failed to verify migrations: %w", err)
	}
	log.Printf("Found %d tables in database", tableCount)

	// List tables
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		return fmt.Errorf("failed to list tables: %w", err)
	}
	defer rows.Close()

	log.Println("Created tables:")
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return fmt.Errorf("failed to scan table name: %w", err)
		}
		log.Printf("- %s", tableName)
	}

	return nil
}

// SeedData seeds the database with initial data
func SeedData(db *sql.DB) error {
	log.Println("Seeding database with initial data...")

	// Read seed SQL file
	seedSQL, err := os.ReadFile("src/database/seed/seed_data.sql")
	if err != nil {
		return fmt.Errorf("failed to read seed file: %w", err)
	}

	// Execute seed SQL
	_, err = db.Exec(string(seedSQL))
	if err != nil {
		return fmt.Errorf("failed to execute seed data: %w", err)
	}

	// Verify seed data
	var teamCount int
	err = db.QueryRow("SELECT COUNT(*) FROM teams").Scan(&teamCount)
	if err != nil {
		return fmt.Errorf("failed to verify seed data: %w", err)
	}
	log.Printf("Seeded %d teams", teamCount)

	var matchCount int
	err = db.QueryRow("SELECT COUNT(*) FROM matches").Scan(&matchCount)
	if err != nil {
		return fmt.Errorf("failed to verify seed data: %w", err)
	}
	log.Printf("Seeded %d matches", matchCount)

	return nil
}
