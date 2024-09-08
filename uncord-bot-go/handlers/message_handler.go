package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

type Handler struct {
	Client   bot.Client
	Lavalink disgolink.Client
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) OnEvent(event bot.Event) {
	switch e := event.(type) {
	case *events.MessageCreate:
		h.OnMessageCreate(e)
	case *events.GuildVoiceStateUpdate:
		h.OnVoiceStateUpdate(e)
	case *events.VoiceServerUpdate:
		h.OnVoiceServerUpdate(e)
	}
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
			_, err := h.Client.Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().
				SetContent("You must be in a voice channel to play music.").Build())
			if err != nil {
				slog.Error("Failed to send message", slog.Any("err", err))
			}
			return
		}

		err := h.play(*guildID, *voiceState.ChannelID, content)
		if err != nil {
			slog.Error("Failed to play track", slog.Any("err", err))
			_, sendErr := h.Client.Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().
				SetContent(fmt.Sprintf("Failed to play the track: %v", err)).Build())
			if sendErr != nil {
				slog.Error("Failed to send error message", slog.Any("err", sendErr))
			}
			return
		}

		_, sendErr := h.Client.Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().
			SetContent("Now playing: "+content).Build())
		if sendErr != nil {
			slog.Error("Failed to send message", slog.Any("err", sendErr))
		}
	}
}

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

func (h *Handler) play(guildID, channelID snowflake.ID, url string) error {
	err := h.Client.UpdateVoiceState(context.Background(), guildID, &channelID, false, false)
	if err != nil {
		return fmt.Errorf("failed to join voice channel: %w", err)
	}

	player := h.Lavalink.Player(guildID)
	if player == nil {
		return fmt.Errorf("could not create player")
	}

	var loadError error
	var trackLoaded bool

	h.Lavalink.BestNode().LoadTracksHandler(context.TODO(), url, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			slog.Info("Loaded a single track", "track", track.Info.Title)
			err := player.Update(context.TODO(), lavalink.WithTrack(track))
			if err != nil {
				loadError = fmt.Errorf("error updating player: %w", err)
			} else {
				trackLoaded = true
			}
		},
		func(playlist lavalink.Playlist) {
			slog.Info("Loaded a playlist", "name", playlist.Info.Name, "trackCount", len(playlist.Tracks))
			if len(playlist.Tracks) > 0 {
				err := player.Update(context.TODO(), lavalink.WithTrack(playlist.Tracks[0]))
				if err != nil {
					loadError = fmt.Errorf("error updating player with playlist: %w", err)
				} else {
					trackLoaded = true
				}
			}
		},
		func(tracks []lavalink.Track) {
			slog.Info("Loaded search results", "trackCount", len(tracks))
			if len(tracks) > 0 {
				err := player.Update(context.TODO(), lavalink.WithTrack(tracks[0]))
				if err != nil {
					loadError = fmt.Errorf("error updating player with search result: %w", err)
				} else {
					trackLoaded = true
				}
			}
		},
		func() {
			loadError = fmt.Errorf("no matches found for URL: %s", url)
		},
		func(err error) {
			loadError = fmt.Errorf("error loading track: %w", err)
		},
	))

	if loadError != nil {
		return loadError
	}

	if !trackLoaded {
		return fmt.Errorf("no track loaded for URL: %s", url)
	}

	slog.Info("Track loaded and playback started", "guildID", guildID, "url", url)
	return nil
}
