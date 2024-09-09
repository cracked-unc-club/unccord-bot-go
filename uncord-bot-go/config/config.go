package config

// Configuration handler for the bot

import (
	"log"
	"os"
	"strconv"

	"github.com/disgoorg/snowflake/v2"
)

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

// Global variable to accesss the configuration
var AppConfig Config

// INitialize the configuration
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

func validateConfig() {
	if AppConfig.DBHost == "" || AppConfig.DBPort == "" || AppConfig.DBUser == "" ||
		AppConfig.DBPassword == "" || AppConfig.DBName == "" || AppConfig.DiscordToken == "" || AppConfig.StarboardChannelID == 0 {
		log.Fatal("One or more required environment variables are missing")
	}
}

// GetEnv returns the value of an environment Variable or a fallback
func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}