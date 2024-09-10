package config

import (
	"log"
	"os"
	"testing"
)

// TestConnectDB tests the database connection with environment variables, only if running in a CI pipeline.
func TestConnectDB(t *testing.T) {
	// Check if we are running in a CI environment (by checking for a 'CI' environment variable)
	if os.Getenv("CI") == "" {
		t.Skip("Skipping TestConnectDB as it's not running in CI environment")
	}

	// Capture original environment variables (using envVars type from the same package).
	originalVariables := &envVars{
		DBHost:             os.Getenv("DB_HOST"),
		DBPort:             os.Getenv("DB_PORT"),
		DBUser:             os.Getenv("DB_USER"),
		DBPassword:         os.Getenv("DB_PASSWORD"),
		DBName:             os.Getenv("DB_NAME"),
	}
	defer originalVariables.setEnvVars() // Restore original environment variables after the test

	// Set the test environment variables.
	testVariables := &envVars{
		DBHost:    "localhost",
		DBPort:    "5432",
		DBUser:    "testuser",
		DBPassword: "testpassword",
		DBName:    "testdb",
	}
	testVariables.setEnvVars()

	// Attempt to connect to the database.
	defer func() {
		if DB != nil {
			DB.Close()
		}
	}()

	// Call the function being tested.
	ConnectDB()

	// Check if DB is not nil (i.e., connection is established).
	if DB == nil {
		t.Fatalf("Expected DB to be initialized, but got nil")
	}

	// Check if the database can be pinged.
	err := DB.Ping()
	if err != nil {
		t.Fatalf("Expected to ping the database successfully, but got error: %v", err)
	}

	log.Println("Test database connection established successfully")
}