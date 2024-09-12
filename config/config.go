package config

// Configuration handler for the bot

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/disgoorg/snowflake/v2"
)

// Config holds the configuration details for the bot, including database credentials, starboard settings, and Discord token.
type Config struct {
	// Database configuration
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string

	// Starboard configuration
	StarboardChannelID snowflake.ID
	StarThreshold      int

	// Other Settings (e.g., Discord Token)
	DiscordToken      string
}

// AppConfig holds the global configuration for the bot.
var AppConfig Config

func LoadConfig() {
    // Parse StarboardChannelID as snowflake.ID
    starboardChannelID, err := snowflake.Parse(os.Getenv("STARBOARD_CHANNEL_ID"))
    if err != nil {
        log.Fatalf("Invalid STARBOARD_CHANNEL_ID: %v", err)
    }

    // Parse StarThreshold as int
    starThreshold, err := strconv.Atoi(os.Getenv("STAR_THRESHOLD"))
    if err != nil {
        log.Fatalf("Invalid STAR_THRESHOLD: %v", err)
    }

    AppConfig = Config{
        DBHost:            os.Getenv("DB_HOST"),
        DBPort:            os.Getenv("DB_PORT"),
        DBUser:            os.Getenv("DB_USER"),
        DBPassword:        os.Getenv("DB_PASSWORD"),
        DBName:            os.Getenv("DB_NAME"),
        StarboardChannelID: starboardChannelID,
        StarThreshold:      starThreshold,
        DiscordToken:      os.Getenv("DISCORD_TOKEN"),
    }

    err = ValidateConfig()
    if err != nil {
        log.Fatalf("Error loading configuration: %v", err)
    }
}

// ValidateConfig checks if the required environment variables are set.
func ValidateConfig() error {
	if AppConfig.DBHost == "" || AppConfig.DBPort == "" || AppConfig.DBUser == "" ||
		AppConfig.DBPassword == "" || AppConfig.DBName == "" || AppConfig.DiscordToken == "" || AppConfig.StarboardChannelID == 0 {
		return fmt.Errorf("one or more required environment variables are missing")
	}
	return nil
}
