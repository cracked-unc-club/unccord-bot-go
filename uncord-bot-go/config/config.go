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
	DBHost    		string
	DBPort  	 	string
	DBUser    		string
	DBPassword 		string
	DBName   		string

	//Starboard configuration
	StarboardChannelID snowflake.ID
	StarThreshold	   int

	// Other Settings (e.g., Discord Token)
	DiscordToken string

}

// AppConfig holds the global configuration for the bot.
var AppConfig Config

// LoadConfig initializes the configuration by loading environment variables and validating them.
func LoadConfig() {
	AppConfig = Config{
		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
		DBUser: os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName: os.Getenv("DB_NAME"),
		StarboardChannelID: snowflake.GetEnv("STARBOARD_CHANNEL_ID"),
		DiscordToken: os.Getenv("DISCORD_TOKEN"),
	}

	// Load and parse star threshold
	starThreshold, err := strconv.Atoi(os.Getenv("STAR_THRESHOLD"))
	if err != nil {
		log.Printf("Error parsing STAR_THRESHOLD: %v", err)
	}
	AppConfig.StarThreshold = starThreshold

	//Validate the configuration to ensure nothing is missing
	validateConfig()

}

func validateConfig() error {
	if AppConfig.DBHost == "" || AppConfig.DBPort == "" || AppConfig.DBUser == "" ||
		AppConfig.DBPassword == "" || AppConfig.DBName == "" || AppConfig.DiscordToken == "" || AppConfig.StarboardChannelID == 0 {
		return fmt.Errorf("one or more required environment variables are missing")
	}
	return nil
}

// // GetEnv returns the value of an environment Variable or a fallback
// func GetEnv(key, fallback string) string {
// 	value, exists := os.LookupEnv(key)
// 	if !exists {
// 		return fallback
// 	}
// 	return value
// }