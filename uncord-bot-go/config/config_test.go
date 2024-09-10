package config

import (
	"bytes"
	"log"
	"os"
	"testing"
)

// envVars is a struct that holds the environment variables
type envVars struct {
    DBHost             string
    DBPort             string
    DBUser             string
    DBPassword         string
    DBName             string
    StarboardChannelID string
    DiscordToken       string
}

// Sets the environment variables based on the struct values
func (e *envVars) setEnvVars() {
    os.Setenv("DB_HOST", e.DBHost)
    os.Setenv("DB_PORT", e.DBPort)
    os.Setenv("DB_USER", e.DBUser)
    os.Setenv("DB_PASSWORD", e.DBPassword)
    os.Setenv("DB_NAME", e.DBName)
    os.Setenv("STARBOARD_CHANNEL_ID", e.StarboardChannelID)
    os.Setenv("DISCORD_TOKEN", e.DiscordToken)
}

// TestEnvironmentVariables tests setting and getting environment variables
func TestEnvironmentVariables(t *testing.T) {
    // Capture original environment variables
    originalVariables := &envVars{
        DBHost:             os.Getenv("DB_HOST"),
        DBPort:             os.Getenv("DB_PORT"),
        DBUser:             os.Getenv("DB_USER"),
        DBPassword:         os.Getenv("DB_PASSWORD"),
        DBName:             os.Getenv("DB_NAME"),
        StarboardChannelID: os.Getenv("STARBOARD_CHANNEL_ID"),
        DiscordToken:       os.Getenv("DISCORD_TOKEN"),
    }
    defer originalVariables.setEnvVars() // Restore original environment variables after the test

    // Define dummy environment variables for testing
    dummyEnvironmentVariables := &envVars{
        DBHost:             "dummyHost",
        DBPort:             "dummyPort",
        DBUser:             "dummyUser",
        DBPassword:         "dummyPassword",
        DBName:             "dummyDBName",
        StarboardChannelID: "dummyChannelID",
        DiscordToken:       "dummyToken",
    }

    // Set dummy environment variables
    dummyEnvironmentVariables.setEnvVars()

    // Define expected environment variables for assertion
    expectedVariables := map[string]string{
        "DB_HOST":             dummyEnvironmentVariables.DBHost,
        "DB_PORT":             dummyEnvironmentVariables.DBPort,
        "DB_USER":             dummyEnvironmentVariables.DBUser,
        "DB_PASSWORD":         dummyEnvironmentVariables.DBPassword,
        "DB_NAME":             dummyEnvironmentVariables.DBName,
        "STARBOARD_CHANNEL_ID": dummyEnvironmentVariables.StarboardChannelID,
        "DISCORD_TOKEN":       dummyEnvironmentVariables.DiscordToken,
    }

    // Assert that the environment variables are as expected
    for key, expected := range expectedVariables {
        if got := os.Getenv(key); got != expected {
            t.Errorf("%s: got %v, want %v", key, got, expected)
        }
    }
}

// TestValidateConfig tests the validateConfig function
func TestValidateConfig(t *testing.T) {

	// Backup original logger and redirect log output to a buffer
    originalLogger := log.Default()
    var logBuffer bytes.Buffer
    log.SetOutput(&logBuffer)
    defer log.SetOutput(originalLogger.Writer()) // Restore original logger after the test


	// Capture original environment variables
    originalVariables := &envVars{
        DBHost:             os.Getenv("DB_HOST"),
        DBPort:             os.Getenv("DB_PORT"),
        DBUser:             os.Getenv("DB_USER"),
        DBPassword:         os.Getenv("DB_PASSWORD"),
        DBName:             os.Getenv("DB_NAME"),
        StarboardChannelID: os.Getenv("STARBOARD_CHANNEL_ID"),
        DiscordToken:       os.Getenv("DISCORD_TOKEN"),
    }
    defer originalVariables.setEnvVars() // Restore original environment variables after the test

	// Define dummy environment variables for testing
	dummyEnvironmentVariables := &envVars{
		DBHost:             "",
		DBPort:             "dummyPort",
		DBUser:             "dummyUser",
		DBPassword:         "dummyPassword",
		DBName:             "dummyDBName",
		StarboardChannelID: "dummyChannelID",
		DiscordToken:       "dummyToken",
	}

	// Set dummy environment variables
	dummyEnvironmentVariables.setEnvVars()

	// Validate the configuration
	err := ValidateConfig()

	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Restore original environment variables after the test
    originalVariables.setEnvVars()
}