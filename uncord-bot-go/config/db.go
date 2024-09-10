package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	// Import the PostgreSQL driver as a blank import, required for the database/sql package.
	_ "github.com/lib/pq"
)

// DB holds the global connection pool to the PostgreSQL database.
var DB *sql.DB

// ConnectDB initializes the database connection using environment variables and establishes a connection pool.
func ConnectDB() {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Cannot ping the database: %v",err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Cannot ping the database: %v", err)
	}
	log.Println("Connected to the database")
}