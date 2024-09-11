package handlers

import (
	"log"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// OnMessageCreate handles message creation events.
func OnMessageCreate(event *events.MessageCreate) {
	if event.Message.Author.Bot {
		return
	}

	// Handle ping-pong responses
	var message string
	switch event.Message.Content {
	case "ping":
		message = "pong"
	case "pong":
		message = "ping"
	default:
	}

	// Send the response back to the channel
	if message != "" {
		_, err := event.Client().Rest().CreateMessage(event.Message.ChannelID, discord.NewMessageCreateBuilder().SetContent(message).Build())
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}
}