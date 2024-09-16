package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"unccord-bot-go/config"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

// InsertStarredMessage handles the insertion of a new starred message into the PostgreSQL database.
func InsertStarredMessage(messageID, channelID, authorID, content string) error {
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
	if !isStarEmoji(event.Emoji) {
		return
	}

	message, err := fetchMessage(event.Client(), event.ChannelID, event.MessageID)
	if err != nil {
		log.Printf("Error fetching message: %v", err)
		return
	}

	if err := updateStarCount(event.MessageID.String(), event.ChannelID.String(), message.Author.ID.String(), message.Content, true); err != nil {
		log.Printf("Error updating star count: %v", err)
		return
	}

	starCount, err := GetStarredMessage(event.MessageID.String())
	if err != nil {
		handleStarCountError(err, event.MessageID.String())
		return
	}

	if err := handleStarboardPost(event, message, starCount); err != nil {
		log.Printf("Error handling starboard post: %v", err)
	}
}

// OnReactionRemove handles the removal of reactions and updates the starboard accordingly.
func OnReactionRemove(event *events.GuildMessageReactionRemove) {
	if !isStarEmoji(event.Emoji) {
		return
	}

	if err := updateStarCount(event.MessageID.String(), event.ChannelID.String(), "", "", false); err != nil {
		log.Printf("Error updating star count: %v", err)
		return
	}

	starCount, err := GetStarredMessage(event.MessageID.String())
	if err != nil {
		handleStarCountError(err, event.MessageID.String())
		return
	}

	if err := updateStarboardMessage(event.Client(), event.MessageID.String(), starCount); err != nil {
		log.Printf("Error updating starboard message: %v", err)
	}
}

// Helper functions

func isStarEmoji(emoji discord.PartialEmoji) bool {
	return emoji.Name != nil && *emoji.Name == "⭐"
}

func fetchMessage(client bot.Client, channelID, messageID snowflake.ID) (*discord.Message, error) {
	message, err := client.Rest().GetMessage(channelID, messageID)
	if err != nil {
		return nil, fmt.Errorf("error fetching message: %w", err)
	}
	return message, nil
}

func updateStarCount(messageID, channelID, authorID, content string, increment bool) error {
	var err error
	if increment {
		err = InsertStarredMessage(messageID, channelID, authorID, content)
	} else {
		err = RemoveStarFromMessage(messageID)
	}
	return err
}

func handleStarCountError(err error, messageID string) {
	if err == sql.ErrNoRows {
		log.Printf("No stars found for message %s, skipping starboard update", messageID)
	} else {
		log.Printf("Error fetching star count: %v", err)
	}
}

func handleStarboardPost(event *events.GuildMessageReactionAdd, message *discord.Message, starCount int) error {
	existingStarboardMessageID, err := GetStarboardMessageID(event.MessageID.String())
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking existing starboard message: %w", err)
	}

	if existingStarboardMessageID != "" {
		return updateStarboardMessage(event.Client(), event.MessageID.String(), starCount)
	}

	if starCount >= config.AppConfig.StarThreshold {
		return PostToStarboard(event, message, starCount)
	}

	return nil
}

func updateStarboardMessage(client bot.Client, messageID string, starCount int) error {
	starboardMessageID, err := GetStarboardMessageID(messageID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil // Message not in starboard yet, nothing to update
		}
		return fmt.Errorf("error fetching starboard message ID: %w", err)
	}

	starboardMessageIDSnowflake, err := snowflake.Parse(starboardMessageID)
	if err != nil {
		return fmt.Errorf("error parsing starboard message ID: %w", err)
	}

	message, err := client.Rest().GetMessage(config.AppConfig.StarboardChannelID, starboardMessageIDSnowflake)
	if err != nil {
		return fmt.Errorf("error fetching starboard message: %w", err)
	}

	if len(message.Embeds) == 0 {
		return fmt.Errorf("starboard message has no embeds")
	}

	updatedEmbed := message.Embeds[0]
	updatedEmbed.Title = fmt.Sprintf("⭐ %d # %s", starCount, message.ChannelID)

	_, err = client.Rest().UpdateMessage(config.AppConfig.StarboardChannelID, starboardMessageIDSnowflake, discord.NewMessageUpdateBuilder().SetEmbeds(updatedEmbed).Build())
	if err != nil {
		return fmt.Errorf("error updating starboard message: %w", err)
	}

	return nil
}

// GetStarboardMessageID retrieves the starboard message ID from the database based on the original message ID.
func GetStarboardMessageID(messageID string) (string, error) {
	var starboardMessageID sql.NullString
	query := `SELECT starboard_message_id FROM starboard WHERE message_id = $1`
	err := config.DB.QueryRow(query, messageID).Scan(&starboardMessageID)
	if err != nil {
		return "", err
	}
	if !starboardMessageID.Valid {
		return "", sql.ErrNoRows
	}
	return starboardMessageID.String, nil
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
		SetTitle(fmt.Sprintf("⭐ %d # %s", starCount, event.ChannelID.String())).                                                                                              // Add star count in title
		SetDescription(message.Content).                                                                                                                                      // Add message content
		AddField("Source", fmt.Sprintf("[Jump!](https://discord.com/channels/%s/%s/%s)", event.GuildID.String(), event.ChannelID.String(), event.MessageID.String()), false). // Jump link to the message
		SetAuthorName(message.Author.Username).                                                                                                                               // Add the author's username
		SetAuthorIcon(avatarURL).                                                                                                                                             // Add the author's avatar URL
		SetTimestamp(message.CreatedAt).                                                                                                                                      // Timestamp of the original message
		SetFooterText("From #" + event.ChannelID.String())                                                                                                                    // Channel name in the footer

	if len(message.Attachments) > 0 {
		embedBuilder.SetImage(message.Attachments[0].URL) // Add the first attachment as an image
	}

	embed := embedBuilder.Build()

	// Send the embed to the starboard channel and capture the message ID
	starboardMessage, err := event.Client().Rest().CreateMessage(config.AppConfig.StarboardChannelID, discord.NewMessageCreateBuilder().AddEmbeds(embed).Build())
	if err != nil {
		return fmt.Errorf("error sending message to starboard: %w", err)
	}

	// Update the database with the starboard message ID
	err = UpdateStarboardMessageID(event.MessageID.String(), starboardMessage.ID.String())
	if err != nil {
		return fmt.Errorf("error updating starboard message ID in database: %w", err)
	}

	return nil
}
