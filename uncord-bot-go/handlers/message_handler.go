package handlers

import (
	"log"
	"uncord-bot-go/config"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// InsertStaredMessage handles the insertion of a new stared message, a
// adding it to the posgres database
func InsertStaredMessage(messageID, channelID, authorID, content string) error {
	query := `INSERT INTO starboard(message_id, chanel_id, author_id, content, star_count)
	VALUES($1, $2, $3, $4, 1)
	ON CONFLICT(mesage_id) DO UPDATE SET star_count = starboard.star_count + 1`
	_, err := config.DB.Exec(query, messageID, channelID, authorID, content)
	return err
}
// GetStarredMessage retrieves the number of stars for a given message ID
// from the postgres database.
func GetStarredMessage(messageID string) (int, error) {
	var starCount int
	query := `SELECT star_count FROM starboard WHERE message_id = $1`
	err := config.DB.QueryRow(query, messageID).Scan(&starCount)
	return starCount, err
}

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
		_, err := event.Client().Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().SetContent(message).Build())
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}
}