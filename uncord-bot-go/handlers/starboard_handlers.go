package handlers

import (
	"log"
	"strconv"
	"uncord-bot-go/config"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// InsertStaredMessage handles the insertion of a new stared message, a
// adding it to the posgres database
func InsertStaredMessage(messageID, channelID, authorID, content string) error {
	query := `INSERT INTO starboard(message_id, channel_id, author_id, content, star_count)
	VALUES($1, $2, $3, $4, 1)
	ON CONFLICT(message_id) DO UPDATE SET star_count = starboard.star_count + 1`
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

//OnReactionAdd handles star reactions and posts the message to the starboard if it reacehs the threshold
func OnReactionAdd(event *events.GuildMessageReactionAdd) {
	// Check if the reaction is a star emoji
	if *event.Emoji.Name == "⭐" {
		//Fetch the message that was reacted to 
		message, err := event.Client().Rest().GetMessage(event.ChannelID, event.MessageID)
		if err != nil {
			log.Printf("Error fetching message: %v", err)
			return
		}

		// Insert or update the star count in the database
		err  = InsertStaredMessage(event.MessageID.String(), event.ChannelID.String(), message.Author.ID.String(), message.Content)
		if err != nil {
			log.Printf("Error inserting stared message: %v", err)
			return
		}

		// check if the message has reached the star threshold
		starCount, err := GetStarredMessage(event.MessageID.String())
		if err != nil {
			log.Printf("Error fetching star count: %v", err)
			return
		}
		
		if starCount >= config.AppConfig.StarThreshold {
			// Call the PostToStarboard function here
			PostToStarboard(event, message, starCount)
		}
	}
}

func PostToStarboard(event *events.GuildMessageReactionAdd, message *discord.Message, starCount int) {
	embed := discord.NewEmbedBuilder().
		SetTitle("Starred Message").
		SetDescription(message.Content).
		AddField("Stars", strconv.Itoa(starCount), true).
		SetAuthorName(message.Author.Username).
		SetTimestamp(message.CreatedAt).
		SetFooterText("From #" + event.ChannelID.String()).
		Build()
	
	_, err := event.Client().Rest().CreateMessage(config.AppConfig.StarboardChannelID, discord.NewMessageCreateBuilder().AddEmbeds(embed).Build())
	if err != nil {
		log.Printf("Error sending message to starboard: %v", err)
		
	}

}

