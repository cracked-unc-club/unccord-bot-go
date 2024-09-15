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

// Config holds the configuration details for the bot, including database credentials, starboard settings, and Discord token.
type Config struct {
	// Database configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Starboard configuration
	StarboardChannelID snowflake.ID
	StarThreshold      int

	// Other Settings (e.g., Discord Token)
	DiscordToken string
}

// AppConfig holds the global configuration for the bot.
var AppConfig Config

func LoadConfig() {
	log.Println("Starting to load configuration...")

	// Load and validate database configuration
	dbConfig := loadDBConfig()

	// Load and validate Starboard configuration
	starboardConfig := loadStarboardConfig()

	// Load and validate Discord token
	discordToken := loadDiscordToken()

	AppConfig = Config{
		DBHost:             dbConfig.Host,
		DBPort:             dbConfig.Port,
		DBUser:             dbConfig.User,
		DBPassword:         dbConfig.Password,
		DBName:             dbConfig.Name,
		StarboardChannelID: starboardConfig.ChannelID,
		StarThreshold:      starboardConfig.Threshold,
		DiscordToken:       discordToken,
	}

	if err := ValidateConfig(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	log.Println("Configuration loaded successfully")
}

func loadDBConfig() struct {
	Host, Port, User, Password, Name string
} {
	return struct {
		Host, Port, User, Password, Name string
	}{
		Host:     getEnvOrFatal("DB_HOST"),
		Port:     getEnvOrFatal("DB_PORT"),
		User:     getEnvOrFatal("DB_USER"),
		Password: getEnvOrFatal("DB_PASSWORD"),
		Name:     getEnvOrFatal("DB_NAME"),
	}
}

func loadStarboardConfig() struct {
	ChannelID snowflake.ID
	Threshold int
} {
	channelIDStr := getEnvOrFatal("STARBOARD_CHANNEL_ID")
	channelID, err := snowflake.Parse(channelIDStr)
	if err != nil {
		log.Fatalf("Invalid STARBOARD_CHANNEL_ID '%s': %v", channelIDStr, err)
	}

	thresholdStr := getEnvOrFatal("STAR_THRESHOLD")
	threshold, err := strconv.Atoi(thresholdStr)
	if err != nil {
		log.Fatalf("Invalid STAR_THRESHOLD '%s': %v", thresholdStr, err)
	}

	return struct {
		ChannelID snowflake.ID
		Threshold int
	}{
		ChannelID: channelID,
		Threshold: threshold,
	}
}

func loadDiscordToken() string {
	token := getEnvOrFatal("DISCORD_TOKEN")
	token = strings.TrimSpace(token)

	log.Printf("Discord token length: %d", len(token))
	log.Printf("Discord token first 10 characters: %s", token[:10])
	log.Printf("Discord token last 10 characters: %s", token[len(token)-10:])

	// Attempt to decode the base64 part of the token
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		log.Fatalf("Token does not have the expected number of parts (expected 3, got %d)", len(parts))
	}

	_, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		log.Fatalf("Failed to decode base64 part of token: %v", err)
	}

	log.Println("Token passed basic validation checks")

	return token
}

func getEnvOrFatal(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is missing or empty", key)
	}
	return value
}

// ValidateConfig checks if the required environment variables are set.
func ValidateConfig() error {
	if AppConfig.DBHost == "" || AppConfig.DBPort == "" || AppConfig.DBUser == "" ||
		AppConfig.DBPassword == "" || AppConfig.DBName == "" || AppConfig.DiscordToken == "" || AppConfig.StarboardChannelID == 0 {
		return fmt.Errorf("one or more required environment variables are missing")
	}
	return nil
}
