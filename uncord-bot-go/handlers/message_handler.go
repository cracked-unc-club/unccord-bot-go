package handlers

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"uncord-bot-go/lavalink"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

type Handler struct {
	lavalink        *lavalink.Client
	client          bot.Client
	voiceStateCache *VoiceStateCache
}

type VoiceStateCache struct {
	cache map[snowflake.ID]map[snowflake.ID]*discord.VoiceState
	mu    sync.RWMutex
}

func NewVoiceStateCache() *VoiceStateCache {
	return &VoiceStateCache{
		cache: make(map[snowflake.ID]map[snowflake.ID]*discord.VoiceState),
	}
}

func (vsc *VoiceStateCache) Update(vs *discord.VoiceState) {
	vsc.mu.Lock()
	defer vsc.mu.Unlock()
	if _, ok := vsc.cache[vs.GuildID]; !ok {
		vsc.cache[vs.GuildID] = make(map[snowflake.ID]*discord.VoiceState)
	}
	vsc.cache[vs.GuildID][vs.UserID] = vs
}

func (vsc *VoiceStateCache) Get(guildID, userID snowflake.ID) (*discord.VoiceState, bool) {
	vsc.mu.RLock()
	defer vsc.mu.RUnlock()
	if guildStates, ok := vsc.cache[guildID]; ok {
		if vs, ok := guildStates[userID]; ok {
			return vs, true
		}
	}
	return nil, false
}

func NewHandler(lavalink *lavalink.Client, client bot.Client) *Handler {
	return &Handler{
		lavalink:        lavalink,
		client:          client,
		voiceStateCache: NewVoiceStateCache(),
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

		voiceState, ok := h.voiceStateCache.Get(*guildID, event.Message.Author.ID)
		if !ok || voiceState.ChannelID == nil {
			_, err := h.client.Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().
				SetContent("You must be in a voice channel to play music.").Build())
			if err != nil {
				slog.Error("Failed to send message", slog.Any("err", err))
			}
			return
		}

		// Play the track
		go func() {
			err := h.lavalink.PlayTrack(*guildID, *voiceState.ChannelID, content)
			if err != nil {
				slog.Error("Failed to play track", slog.Any("err", err))
				_, sendErr := h.client.Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().
					SetContent(fmt.Sprintf("Failed to play the track: %v", err)).Build())
				if sendErr != nil {
					slog.Error("Failed to send error message", slog.Any("err", sendErr))
				}
				return
			}

			_, err = h.client.Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().
				SetContent("Now playing: "+content).Build())
			if err != nil {
				slog.Error("Failed to send message", slog.Any("err", err))
			}
		}()
	}
}

func (h *Handler) handleVoiceStateUpdate(event *events.GuildVoiceStateUpdate) {
	h.voiceStateCache.Update(&event.VoiceState)
	if event.VoiceState.UserID == h.client.ApplicationID() {
		slog.Info("Voice state updated", "guildID", event.VoiceState.GuildID, "channelID", event.VoiceState.ChannelID, "sessionID", event.VoiceState.SessionID)
		h.lavalink.OnVoiceStateUpdate(event.VoiceState.GuildID, *event.VoiceState.ChannelID, event.VoiceState.SessionID)
	}
}

func (h *Handler) handleVoiceServerUpdate(event *events.VoiceServerUpdate) {
	slog.Info("Voice server updated", "guildID", event.GuildID, "endpoint", *event.Endpoint, "token", event.Token)
	h.lavalink.OnVoiceServerUpdate(event.GuildID, event.Token, *event.Endpoint)
}
