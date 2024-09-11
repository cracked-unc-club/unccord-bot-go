package internal

// Initial utility funciton LogError is useless (but an example)
// Add more useful functions as needed

import (
	"log"
)

// LogError logs an error if one is present.
// This function is a utility example that can be extended as needed.
func LogError(err error) {
	if err != nil {
		log.Printf("Error: %v", err)
	}
}