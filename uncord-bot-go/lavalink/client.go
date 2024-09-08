package lavalink

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

const UNC_DEC = 11711099
const lavalinkUserId = snowflake.ID(UNC_DEC)

type Client struct {
	link             disgolink.Client
	connectionStates map[snowflake.ID]*connectionState
	mu               sync.Mutex
	discordClient    bot.Client
}

type connectionState struct {
	connected  bool
	readyChan  chan struct{}
	retryCount int
}

func NewClient(nodeConfig disgolink.NodeConfig, discordClient bot.Client) (*Client, error) {
	client := &Client{
		connectionStates: make(map[snowflake.ID]*connectionState),
		discordClient:    discordClient,
	}

	link := disgolink.New(lavalinkUserId,
		disgolink.WithListenerFunc(func(player disgolink.Player, event lavalink.WebSocketClosedEvent) {
			slog.Info("WebSocket closed", "guildID", event.GuildID(), "code", event.Code, "reason", event.Reason)
			client.handleDisconnect(event.GuildID())
		}),
		disgolink.WithListenerFunc(func(player disgolink.Player, event lavalink.TrackStuckEvent) {
			slog.Info("Track stuck", "guildID", event.GuildID(), "track", event.Track.Info.Title)
		}),
		disgolink.WithListenerFunc(func(player disgolink.Player, event lavalink.TrackExceptionEvent) {
			slog.Info("Track exception", "guildID", event.GuildID(), "track", event.Track.Info.Title, "error", event.Exception)
		}),
	)
	_, err := link.AddNode(context.TODO(), nodeConfig)
	if err != nil {
		return nil, err
	}
	client.link = link
	return client, nil
}

func (c *Client) PlayTrack(guildID snowflake.ID, channelID snowflake.ID, url string) error {
	// First, ensure we're connected to the voice channel
	err := c.ensureConnected(guildID, channelID)
	if err != nil {
		return fmt.Errorf("failed to connect to voice channel: %w", err)
	}

	// Wait for the connection to be established
	time.Sleep(5 * time.Second)

	player := c.link.Player(guildID)
	var loadError error
	var trackLoaded bool

	c.link.BestNode().LoadTracksHandler(context.TODO(), url, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			slog.Info("Loaded a single track", "track", track.Info.Title)
			err := player.Update(context.TODO(), lavalink.WithTrack(track), lavalink.WithPaused(false))
			if err != nil {
				loadError = fmt.Errorf("error updating player: %w", err)
			} else {
				trackLoaded = true
			}
		},
		func(playlist lavalink.Playlist) {
			slog.Info("Loaded a playlist", "name", playlist.Info.Name, "trackCount", len(playlist.Tracks))
			if len(playlist.Tracks) > 0 {
				err := player.Update(context.TODO(), lavalink.WithTrack(playlist.Tracks[0]), lavalink.WithPaused(false))
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
				err := player.Update(context.TODO(), lavalink.WithTrack(tracks[0]), lavalink.WithPaused(false))
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

func (c *Client) OnVoiceStateUpdate(guildID, channelID snowflake.ID, sessionID string) {
	c.link.OnVoiceStateUpdate(context.TODO(), guildID, &channelID, sessionID)
	c.setConnected(guildID)
}

func (c *Client) OnVoiceServerUpdate(guildID snowflake.ID, token string, endpoint string) {
	c.link.OnVoiceServerUpdate(context.TODO(), guildID, token, endpoint)
	c.setConnected(guildID)
}

func (c *Client) setConnected(guildID snowflake.ID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	state, ok := c.connectionStates[guildID]
	if !ok {
		state = &connectionState{readyChan: make(chan struct{})}
		c.connectionStates[guildID] = state
	}
	if !state.connected {
		state.connected = true
		close(state.readyChan)
	}
}

func (c *Client) ensureConnected(guildID, channelID snowflake.ID) error {
	c.mu.Lock()
	state, ok := c.connectionStates[guildID]
	if !ok {
		state = &connectionState{readyChan: make(chan struct{})}
		c.connectionStates[guildID] = state
	}
	c.mu.Unlock()

	if !state.connected {
		err := c.discordClient.UpdateVoiceState(context.Background(), guildID, &channelID, false, false)
		if err != nil {
			return fmt.Errorf("failed to join voice channel: %w", err)
		}

		return nil
	}

	return nil
}

func (c *Client) handleDisconnect(guildID snowflake.ID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	state, ok := c.connectionStates[guildID]
	if !ok {
		return
	}

	state.connected = false
	state.retryCount++

	if state.retryCount > 3 {
		delete(c.connectionStates, guildID)
		slog.Info("Removed connection state after multiple reconnection attempts", "guildID", guildID)
		return
	}

	state.readyChan = make(chan struct{})

	go func() {
		// Implement reconnection logic here if needed
	}()
}
