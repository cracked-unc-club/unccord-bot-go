package config

// Configuration handler for the bot

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/disgoorg/snowflake/v2"
)

// Config holds the configuration details for the bot, including database credentials, starboard settings, Discord token, and Lavalink configuration.
type Config struct {
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	StarboardChannelID snowflake.ID
	StarThreshold     int
	LavalinkHost      string
	LavalinkPort      string
	LavalinkPassword  string
	DiscordToken      string
}

// AppConfig holds the global configuration for the bot.
var AppConfig Config

// LoadConfig initializes and validates the bot's configuration.
func LoadConfig() {
	log.Println("Starting to load configuration...")

	AppConfig = Config{
		DBHost:             mustGetEnv("DB_HOST"),
		DBPort:             mustGetEnv("DB_PORT"),
		DBUser:             mustGetEnv("DB_USER"),
		DBPassword:         mustGetEnv("DB_PASSWORD"),
		DBName:             mustGetEnv("DB_NAME"),
		StarboardChannelID: mustParseSnowflake("STARBOARD_CHANNEL_ID"),
		StarThreshold:      mustParseInt("STAR_THRESHOLD"),
		DiscordToken:       loadAndValidateDiscordToken(),
		LavalinkHost:       mustGetEnv("SERVER_ADDRESS"),
		LavalinkPort:       mustGetEnv("SERVER_PORT"),
		LavalinkPassword:   mustGetEnv("LAVALINK_SERVER_PASSWORD"),
	}

	if err := ValidateConfig(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	log.Println("Configuration loaded successfully")
}

// mustGetEnv retrieves environment variables or logs a fatal error if not set.
func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is missing or empty", key)
	}
	return value
}

// mustParseSnowflake parses a snowflake ID from an environment variable.
func mustParseSnowflake(key string) snowflake.ID {
	idStr := mustGetEnv(key)
	id, err := snowflake.Parse(idStr)
	if err != nil {
		log.Fatalf("Invalid %s: %v", key, err)
	}
	return id
}

// mustParseInt parses an integer from an environment variable.
func mustParseInt(key string) int {
	valueStr := mustGetEnv(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Fatalf("Invalid %s: %v", key, err)
	}
	return value
}

// loadAndValidateDiscordToken handles Discord token retrieval and basic validation.
func loadAndValidateDiscordToken() string {
	token := strings.TrimSpace(mustGetEnv("DISCORD_TOKEN"))

	log.Printf("Discord token length: %d", len(token))
	log.Printf("Discord token first 10 characters: %s", token[:10])
	log.Printf("Discord token last 10 characters: %s", token[len(token)-10:])

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		log.Fatalf("Token does not have the expected number of parts (expected 3, got %d)", len(parts))
	}

	if _, err := base64.RawStdEncoding.DecodeString(parts[0]); err != nil {
		log.Fatalf("Failed to decode base64 part of token: %v", err)
	}

	log.Println("Token passed basic validation checks")
	return token
}

// ValidateConfig checks if all critical configuration fields are set.
func ValidateConfig() error {
	if AppConfig.DBHost == "" || AppConfig.DBPort == "" || AppConfig.DBUser == "" ||
		AppConfig.DBPassword == "" || AppConfig.DBName == "" || AppConfig.DiscordToken == "" ||
		AppConfig.LavalinkHost == "" || AppConfig.LavalinkPort == "" || AppConfig.LavalinkPassword == "" {
		return fmt.Errorf("one or more required environment variables are missing")
	}
	return nil
}
