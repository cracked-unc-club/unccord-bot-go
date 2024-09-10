package config

import (
	"bytes"
	"log"
	"os"
	"testing"
)

// TestConnectDB simulates a failure to connect to the database using incorrect environment variables.
func TestConnectDB(t *testing.T) {
	// Skip test if not in a CI environment
	if os.Getenv("CI") == "" {
		t.Skip("Skipping TestConnectDB as it's not running in CI environment")
	}

	// Capture the log output
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr) // Restore original log output after the test

	// Set incorrect environment variables to force a failure
	os.Setenv("DB_HOST", "invalid_host")
	os.Setenv("DB_PORT", "invalid_port")
	os.Setenv("DB_USER", "invalid_user")
	os.Setenv("DB_PASSWORD", "invalid_password")
	os.Setenv("DB_NAME", "invalid_db")

	// Call the ConnectDB function
	ConnectDB()

	// Check the log output for the expected error message
	if !bytes.Contains(logBuffer.Bytes(), []byte("Cannot ping the database")) {
		t.Errorf("Expected error message not found in log output: %v", logBuffer.String())
	}
}