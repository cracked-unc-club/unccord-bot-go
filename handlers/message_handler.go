package handlers

import (
	"log"
	"github.com/disgoorg/disgo/events"
)

func MessageHandler(event *event.MessageCreate) {
	if event.Message.Author.Bot {
		return
	}

	switch event.Message.Content {
	case "!ping":
		_, err :- event.Client().Rest().CreateMessage(event.ChannelID, "Pong!")
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