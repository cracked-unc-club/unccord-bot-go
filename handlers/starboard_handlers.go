package handlers

import (
	"fmt"
	"log"
	"uncord-bot-go/config"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

// InsertStaredMessage handles the insertion of a new starred message into the PostgreSQL database.
func InsertStaredMessage(messageID, channelID, authorID, content string) error {
	query := `INSERT INTO starboard(message_id, channel_id, author_id, content, star_count)
	VALUES($1, $2, $3, $4, 1)
	ON CONFLICT(message_id) DO UPDATE SET star_count = starboard.star_count + 1`
	_, err := config.DB.Exec(query, messageID, channelID, authorID, content)
	return err
}

// GetStarredMessage retrieves the number of stars for a given message ID from the PostgreSQL database.
func GetStarredMessage(messageID string) (int, error) {
	var starCount int
	query := `SELECT star_count FROM starboard WHERE message_id = $1`
	err := config.DB.QueryRow(query, messageID).Scan(&starCount)
	return starCount, err
}

// OnReactionAdd handles star reactions and posts the message to the starboard if it reaches the threshold.
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

// OnReactionRemove handles the removal of reactions and updates the starboard accordingly.
func OnReactionRemove(event *events.GuildMessageReactionRemove) {
	// Check if the removed reaction is a star emoji
	if event.Emoji.Name != nil && *event.Emoji.Name == "⭐" {
		// Decrement the star count in the database
		err := RemoveStarFromMessage(event.MessageID.String())
		if err != nil {
			log.Printf("Error removing star from message: %v", err)
			return
		}

		// Fetch the updated star count
		starCount, err := GetStarredMessage(event.MessageID.String())
		if err != nil {
			log.Printf("Error fetching star count: %v", err)
			return
		}

		// If the star count is 0, remove the message from the starboard
		if starCount <= 0 {
			// Fetch the starboard message ID from the database
			starboardMessageID, err := GetStarboardMessageID(event.MessageID.String())
			if err != nil {
				log.Printf("Error fetching starboard message ID: %v", err)
				return
			}

			// Delete the message from the starboard channel
			err = DeleteStarboardMessage(event.Client(), starboardMessageID)
			if err != nil {
				log.Printf("Error deleting message from starboard: %v", err)
			}

			// Optionally, you could delete the row from the database
			// err = RemoveFromStarboard(event.MessageID.String())
			// if err != nil {
			//     log.Printf("Error removing message from starboard: %v", err)
			// }
		}
	}
}

// GetStarboardMessageID retrieves the starboard message ID from the database based on the original message ID.
func GetStarboardMessageID(messageID string) (string, error) {
	var starboardMessageID string
	query := `SELECT starboard_message_id FROM starboard WHERE message_id = $1`
	err := config.DB.QueryRow(query, messageID).Scan(&starboardMessageID)
	return starboardMessageID, err
}

// DeleteStarboardMessage deletes a message from the starboard channel.
func DeleteStarboardMessage(client bot.Client, starboardMessageID string) error {
	// Parse the starboard message ID from string to snowflake.ID
	starboardMessageIDSnowflake, err := snowflake.Parse(starboardMessageID)
	if err != nil {
		log.Printf("Error parsing starboard message ID: %v", err)
		return err
	}

	// Delete the message from the starboard channel
	err = client.Rest().DeleteMessage(config.AppConfig.StarboardChannelID, starboardMessageIDSnowflake)
	if err != nil {
		log.Printf("Error deleting message from starboard: %v", err)
	}
	return err
}

// RemoveStarFromMessage decreases the star count of a message in the PostgreSQL database.
func RemoveStarFromMessage(messageID string) error {
	query := `UPDATE starboard SET star_count = star_count - 1 WHERE message_id = $1 AND star_count > 0`
	_, err := config.DB.Exec(query, messageID)
	return err
}

// RemoveFromStarboard deletes a message from the starboard in the PostgreSQL database.
func RemoveFromStarboard(messageID string) error {
	query := `DELETE FROM starboard WHERE message_id = $1`
	_, err := config.DB.Exec(query, messageID)
	return err
}

// UpdateStarboardMessageID updates the starboard message ID in the PostgreSQL database after the message is posted to the starboard.
func UpdateStarboardMessageID(messageID, starboardMessageID string) error {
	query := `UPDATE starboard SET starboard_message_id = $1 WHERE message_id = $2`
	_, err := config.DB.Exec(query, starboardMessageID, messageID)
	return err
}

// PostToStarboard posts a message to the starboard and updates the database with the starboard message ID.
func PostToStarboard(event *events.GuildMessageReactionAdd, message *discord.Message, starCount int) error {
	// Safely handle the author's avatar URL
	avatarURL := ""
	if message.Author.AvatarURL() != nil {
		avatarURL = *message.Author.AvatarURL()
	}

	// Create the embed for the starred message
	embedBuilder := discord.NewEmbedBuilder().
		SetTitle(fmt.Sprintf("⭐ %d # %s", starCount, event.ChannelID.String())). // Add star count in title
		SetDescription(message.Content).                                          // Add message content
		AddField("Source", fmt.Sprintf("[Jump!](https://discord.com/channels/%s/%s/%s)", event.GuildID.String(), event.ChannelID.String(), event.MessageID.String()), false). // Jump link to the message
		SetAuthorName(message.Author.Username).                                   // Add the author's username
		SetAuthorIcon(avatarURL).                                                 // Add the author's avatar URL
		SetTimestamp(message.CreatedAt).                                          // Timestamp of the original message
		SetFooterText("From #" + event.ChannelID.String())                        // Channel name in the footer

	if len(message.Attachments) > 0 {
		embedBuilder.SetImage(message.Attachments[0].URL) // Add the first attachment as an image
	}

	embed := embedBuilder.Build()

	// Send the embed to the starboard channel and capture the message ID
	starboardMessage, err := event.Client().Rest().CreateMessage(config.AppConfig.StarboardChannelID, discord.NewMessageCreateBuilder().AddEmbeds(embed).Build())
	if err != nil {
		log.Printf("Error sending message to starboard: %v", err)
		return err
	}

	// Update the database with the starboard message ID
	err = UpdateStarboardMessageID(event.MessageID.String(), starboardMessage.ID.String())
	if err != nil {
		log.Printf("Error updating starboard message ID in database: %v", err)
		return err
	}

	return nil
}

