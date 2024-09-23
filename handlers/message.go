package handlers

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

type PlayerInfo struct {
	TrackInfo lavalink.TrackInfo
	IsPlaying bool
}

func (h *Handler) OnMessageCreate(event *events.MessageCreate) {
	if event.Message.Author.Bot {
		return
	}

	content := event.Message.Content
	if strings.HasPrefix(content, "http://") || strings.HasPrefix(content, "https://") {
		guildID := event.Message.GuildID
		if guildID == nil {
			return
		}

		voiceState, exists := h.Client.Caches().VoiceState(*guildID, event.Message.Author.ID)
		if !exists || voiceState.ChannelID == nil {
			return
		}

		err := h.play(*guildID, event.ChannelID, *voiceState.ChannelID, content)
		if err != nil {
			slog.Error("Failed to play track", slog.Any("err", err))
			embed := discord.NewEmbedBuilder().
				SetTitle("Error").
				SetDescription(fmt.Sprintf("Failed to play the track: %v", err)).
				SetColor(ColorError).
				Build()

			_, sendErr := h.Client.Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().
				SetEmbeds(embed).
				SetEphemeral(true).
				Build())
			if sendErr != nil {
				slog.Error("Failed to send error message", slog.Any("err", sendErr))
			}
		}
	}
}
