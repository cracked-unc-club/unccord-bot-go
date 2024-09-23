package handlers

import (
	"context"
	"log/slog"

	"github.com/disgoorg/disgo/events"
)

func (h *Handler) OnVoiceStateUpdate(event *events.GuildVoiceStateUpdate) {
	if event.VoiceState.UserID == h.Client.ApplicationID() {
		slog.Info("Voice state updated", "guildID", event.VoiceState.GuildID, "channelID", event.VoiceState.ChannelID, "sessionID", event.VoiceState.SessionID)
		h.Lavalink.OnVoiceStateUpdate(context.TODO(), event.VoiceState.GuildID, event.VoiceState.ChannelID, event.VoiceState.SessionID)
	}
}

func (h *Handler) OnVoiceServerUpdate(event *events.VoiceServerUpdate) {
	slog.Info("Voice server updated", "guildID", event.GuildID, "endpoint", *event.Endpoint, "token", event.Token)
	h.Lavalink.OnVoiceServerUpdate(context.TODO(), event.GuildID, event.Token, *event.Endpoint)
}
