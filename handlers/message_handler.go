// Package handlers provides functions for handling Discord message events.
// This file contains the MessageHandler function, which handles incoming message events and responds accordingly.
package handlers

import (
	"log"
	"github.com/disgoorg/disgo/events"
)

// MessageHandler is a function that handles incoming message events and responds accordingly.
// It checks the content of the message and performs different actions based on the command.
// If the message content is "!ping", it sends a "Pong!" response.
// If the message content is "!hello", it sends a "Hello there!" response.
// For any other message content, it sends a "Whoops, Unknown command" response.
// The function uses the event.Client().Rest().CreateMessage method to send the responses.
// If there is an error sending the message, it logs the error.
func MessageHandler(event *event.MessageCreate) {
	if event.Message.Author.Bot {
		return
	}

	switch event.Message.Content {
	case "!ping":
		_, err := event.Client().Rest().CreateMessage(event.ChannelID, "Pong!")
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
	case "!hello":
		_, err := event.Client().Rest().CreateMessage(event.ChannelID, "Hello there!")
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
	default:
		// handle unknown commands or other logic
		_, err := event.Client().Rest().CreateMessage(event.ChannelID, "Whoops, Unknown command")
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
}