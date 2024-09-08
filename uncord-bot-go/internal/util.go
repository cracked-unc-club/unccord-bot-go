package internal

// Initial utility funciton LogError is useless (but an example)
// Add more useful functions as needed

import "log"

func LogError(err error) {
	if err != nil {
		log.Printf("Error: %v", err)
	}
}
