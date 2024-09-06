package config
// Configuration handler for the bot

import "os"

// GetEnv returns the value of an environment Variable or a fallback
func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}