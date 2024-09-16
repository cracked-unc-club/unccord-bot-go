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
		log.Fatalf("Cannot open the database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Cannot ping the database: %v", err)
	}
	log.Println("Connected to the database")

	// Initialize the database schema
	err = initDBSchema()
	if err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}
	log.Println("Database schema initialized")
}

// Initialize the database schema
func initDBSchema() error {
	// Check if the starboard table already exists
	var exists bool
	err := DB.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'starboard')").Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if starboard table exists: %v", err)
	}

	if exists {
		log.Println("Starboard table already exists, skipping initialization")
		return nil
	}

	// Read and execute the SQL script only if the table doesn't exist
	script, err := os.ReadFile("SQL/starboard-ddl.sql")
	if err != nil {
		return fmt.Errorf("failed to read starboard-ddl.sql: %v", err)
	}

	_, err = DB.Exec(string(script))
	if err != nil {
		return fmt.Errorf("failed to execute starboard-ddl.sql: %v", err)
	}

	log.Println("Starboard table created successfully")
	return nil
}
