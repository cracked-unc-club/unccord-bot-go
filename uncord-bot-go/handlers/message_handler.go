package handlers

import (
	"fmt"
	"log/slog"
	"strings"
	"uncord-bot-go/lavalink"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type Handler struct {
	lavalink *lavalink.Client
	client   bot.Client
}

func NewHandler(lavalink *lavalink.Client, client bot.Client) *Handler {
	return &Handler{
		lavalink: lavalink,
		client:   client,
	}
}

func (h *Handler) OnEvent(event bot.Event) {
	switch e := event.(type) {
	case *events.MessageCreate:
		h.handleMessageCreate(e)
	case *events.GuildVoiceStateUpdate:
		h.handleVoiceStateUpdate(e)
	case *events.VoiceServerUpdate:
		h.handleVoiceServerUpdate(e)
	}
}

func (h *Handler) handleMessageCreate(event *events.MessageCreate) {
	if event.Message.Author.Bot {
		return
	}

	content := event.Message.Content
	if strings.HasPrefix(content, "http://") || strings.HasPrefix(content, "https://") {
		guildID := event.Message.GuildID
		if guildID == nil {
			return
		}

		voiceState, err := h.client.Rest().GetUserVoiceState(*guildID, event.Message.Author.ID)
		if err != nil || voiceState == nil || voiceState.ChannelID == nil {
			_, err = h.client.Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().
				SetContent("You must be in a voice channel to play music.").Build())
			if err != nil {
				slog.Error("Failed to send message", slog.Any("err", err))
			}
			return
		}

		// Play the track
		err = h.lavalink.PlayTrack(*guildID, *voiceState.ChannelID, content)
		if err != nil {
			slog.Error("Failed to play track", slog.Any("err", err))
			_, err = h.client.Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().
				SetContent(fmt.Sprintf("Failed to play the track: %v", err)).Build())
			if err != nil {
				slog.Error("Failed to send message", slog.Any("err", err))
			}
			return
		}

		_, err = h.client.Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().
			SetContent("Now playing: "+content).Build())
		if err != nil {
			slog.Error("Failed to send message", slog.Any("err", err))
		}
	}
}

func (h *Handler) handleVoiceStateUpdate(event *events.GuildVoiceStateUpdate) {
	if event.VoiceState.UserID == h.client.ApplicationID() {
		slog.Info("Voice state updated", "guildID", event.VoiceState.GuildID, "channelID", event.VoiceState.ChannelID, "sessionID", event.VoiceState.SessionID)
		h.lavalink.OnVoiceStateUpdate(event.VoiceState.GuildID, *event.VoiceState.ChannelID, event.VoiceState.SessionID)
	}
}

func (h *Handler) handleVoiceServerUpdate(event *events.VoiceServerUpdate) {
	slog.Info("Voice server updated", "guildID", event.GuildID, "endpoint", *event.Endpoint, "token", event.Token)
	h.lavalink.OnVoiceServerUpdate(event.GuildID, event.Token, *event.Endpoint)
}
