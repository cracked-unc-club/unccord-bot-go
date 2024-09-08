package handlers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

func (h *Handler) OnComponentInteraction(event *events.ComponentInteractionCreate) {
	if event.GuildID() == nil {
		return
	}

	switch event.Data.CustomID() {
	case "playpause":
		h.handlePlayPause(event)
	case "skip":
		h.handleSkipButton(event)
	case "rewind":
		h.handleRewind(event)
	}
}

func (h *Handler) play(guildID, commandChannelID, voiceChannelID snowflake.ID, url string) error {
	err := h.Client.UpdateVoiceState(context.Background(), guildID, &voiceChannelID, false, false)
	if err != nil {
		return fmt.Errorf("failed to join voice channel: %w", err)
	}

	player := h.Lavalink.Player(guildID)
	queue := h.Queues.Get(guildID)

	var loadError error
	var trackLoaded bool
	var addedTrack lavalink.Track

	h.Lavalink.BestNode().LoadTracksHandler(context.TODO(), url, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			queuePosition := len(queue.Tracks)
			queue.Add(track)
			addedTrack = track

			if queuePosition == 0 {
				err := h.playTrack(guildID, track)
				if err != nil {
					slog.Error("Failed to play track", slog.Any("err", err))
				} else {
					trackLoaded = true
					// Post the player control panel only for the first song
					h.createControlPanel(commandChannelID, guildID)
				}
			} else {
				slog.Info("Added track to queue", "track", track.Info.Title, "position", queuePosition+1)
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

	// Send appropriate message based on queue position
	queuePosition := len(queue.Tracks) - 1
	var embed *discord.EmbedBuilder
	if queuePosition == 0 {
		embed = discord.NewEmbedBuilder().
			SetTitle("Now Playing").
			SetDescription(fmt.Sprintf("**%s**\nby %s", addedTrack.Info.Title, addedTrack.Info.Author)).
			SetColor(ColorSuccess).
			SetThumbnail(*addedTrack.Info.ArtworkURL)
	} else {
		embed = discord.NewEmbedBuilder().
			SetTitle("Added to Queue").
			SetDescription(fmt.Sprintf("**%s**\nby %s\n\nPosition in queue: %d",
				addedTrack.Info.Title,
				addedTrack.Info.Author,
				queuePosition+1)).
			SetColor(ColorInfo).
			SetThumbnail(*addedTrack.Info.ArtworkURL)
	}

	_, err = h.Client.Rest().CreateMessage(commandChannelID, discord.NewMessageCreateBuilder().
		SetEmbeds(embed.Build()).
		Build())
	if err != nil {
		slog.Error("Failed to send queue message", slog.Any("err", err))
	}

	return nil
}

func (h *Handler) sendNowPlayingMessage(channelID snowflake.ID, track lavalink.Track) {
	_, err := h.Client.Rest().CreateMessage(channelID, discord.NewMessageCreateBuilder().
		SetContent(fmt.Sprintf("Now playing: %s by %s", track.Info.Title, track.Info.Author)).
		SetEphemeral(true).
		Build())
	if err != nil {
		slog.Error("Failed to send now playing message", slog.Any("err", err))
	}
}

func (h *Handler) sendAddedToQueueMessage(channelID snowflake.ID, track lavalink.Track, position int) {
	_, err := h.Client.Rest().CreateMessage(channelID, discord.NewMessageCreateBuilder().
		SetContent(fmt.Sprintf("Added to queue: %s by %s (Position: %d)", track.Info.Title, track.Info.Author, position+1)).
		SetEphemeral(true).
		Build())
	if err != nil {
		slog.Error("Failed to send added to queue message", slog.Any("err", err))
	}
}

func (h *Handler) createControlPanel(channelID, guildID snowflake.ID) {
	queue := h.Queues.Get(guildID)
	if len(queue.Tracks) == 0 {
		_, err := h.Client.Rest().CreateMessage(channelID, discord.NewMessageCreateBuilder().
			SetContent("No songs in the queue.").
			SetEphemeral(true).
			Build())
		if err != nil {
			slog.Error("Failed to send empty queue message", slog.Any("err", err))
		}
		return
	}

	currentTrack := queue.Tracks[0]

	_, err := h.Client.Rest().CreateMessage(channelID, discord.NewMessageCreateBuilder().
		SetContent("").
		SetEmbeds(discord.NewEmbedBuilder().
			SetTitle(currentTrack.Info.Title).
			SetDescription(currentTrack.Info.Author).
			SetColor(ColorInfo).
			SetImage(*currentTrack.Info.ArtworkURL).
			Build(),
		).
		AddActionRow(
			discord.NewSecondaryButton("⏪ Rewind", "rewind"),
			discord.NewPrimaryButton("⏯️ Play/Pause", "playpause"),
			discord.NewSecondaryButton("⏩ Skip", "skip"),
		).
		Build(),
	)

	if err != nil {
		slog.Error("Failed to create control panel", slog.Any("err", err))
	}
}

func (h *Handler) handlePlayPause(event *events.ComponentInteractionCreate) {
	player := h.Lavalink.Player(*event.GuildID())
	if player == nil {
		_ = event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("No active player found.").SetEphemeral(true).Build())
		return
	}

	var action string
	if player.Paused() {
		_ = player.Update(context.TODO(), lavalink.WithPaused(false))
		action = "resumed"
	} else {
		_ = player.Update(context.TODO(), lavalink.WithPaused(true))
		action = "paused"
	}

	_ = event.CreateMessage(discord.NewMessageCreateBuilder().SetContent(fmt.Sprintf("Playback %s.", action)).SetEphemeral(true).Build())
}

func (h *Handler) handleRewind(event *events.ComponentInteractionCreate) {
	player := h.Lavalink.Player(*event.GuildID())
	if player == nil {
		_ = event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("No active player found.").SetEphemeral(true).Build())
		return
	}

	currentPosition := player.Position()
	newPosition := currentPosition - lavalink.Duration(10_000*lavalink.Millisecond)
	if newPosition < 0 {
		newPosition = 0
	}

	_ = player.Update(context.TODO(), lavalink.WithPosition(newPosition))
	_ = event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Rewound 10 seconds.").SetEphemeral(true).Build())
}

func (h *Handler) playTrack(guildID snowflake.ID, track lavalink.Track) error {
	player := h.Lavalink.Player(guildID)
	err := player.Update(context.TODO(), lavalink.WithTrack(track))
	if err != nil {
		slog.Error("Error updating player", slog.Any("err", err))
		return err
	}

	player.Lavalink().AddListeners(disgolink.NewListenerFunc(func(player disgolink.Player, event lavalink.Event) {
		switch evt := event.(type) {
		case *lavalink.TrackEndEvent:
			if evt.Reason == lavalink.TrackEndReasonFinished {
				h.playNextTrack(guildID)
			}
		}
	}))

	return nil
}

func (h *Handler) playNextTrack(guildID snowflake.ID) {
	queue := h.Queues.Get(guildID)
	if len(queue.Tracks) == 0 {
		// If there are no more tracks, stop the player
		player := h.Lavalink.ExistingPlayer(guildID)
		if player != nil {
			if err := player.Update(context.TODO(), lavalink.WithNullTrack()); err != nil {
				slog.Error("Failed to stop player", slog.Any("err", err))
			}
		}
		slog.Info("Queue ended, stopped player", "guildID", guildID)
		return
	}

	nextTrack := queue.Tracks[0]
	queue.Tracks = queue.Tracks[1:]

	err := h.playTrack(guildID, nextTrack)
	if err != nil {
		slog.Error("Failed to play next track", slog.Any("err", err))
	}
}

func (h *Handler) handleSkipButton(event *events.ComponentInteractionCreate) {
	embed, err := h.skipTracks(*event.GuildID(), 1)
	if err != nil {
		event.CreateMessage(discord.NewMessageCreateBuilder().
			SetEmbeds(discord.NewEmbedBuilder().
				SetDescription(fmt.Sprintf("Error: %s", err)).
				SetColor(ColorError).
				Build()).
			SetEphemeral(true).
			Build())
		return
	}

	event.CreateMessage(discord.NewMessageCreateBuilder().
		SetEmbeds(embed.Build()).
		Build())
}

func (h *Handler) skipTracks(guildID snowflake.ID, amount int) (*discord.EmbedBuilder, error) {
	player := h.Lavalink.ExistingPlayer(guildID)
	queue := h.Queues.Get(guildID)
	if player == nil {
		return nil, fmt.Errorf("no player found")
	}

	if len(queue.Tracks) == 0 {
		// If queue is empty, stop the current track
		if err := player.Update(context.TODO(), lavalink.WithNullTrack()); err != nil {
			return nil, fmt.Errorf("error while stopping track: %w", err)
		}
		return discord.NewEmbedBuilder().
			SetDescription("No more tracks in the queue. Stopped playing.").
			SetColor(ColorInfo), nil
	}

	skippedTracks := min(amount, len(queue.Tracks))
	queue.Tracks = queue.Tracks[skippedTracks:]

	if len(queue.Tracks) == 0 {
		// If we've skipped all tracks, stop the player
		if err := player.Update(context.TODO(), lavalink.WithNullTrack()); err != nil {
			return nil, fmt.Errorf("error while stopping track: %w", err)
		}
		return discord.NewEmbedBuilder().
			SetDescription("Skipped all tracks. No more tracks in the queue. Stopped playing.").
			SetColor(ColorInfo), nil
	}

	nextTrack := queue.Tracks[0]
	if err := player.Update(context.TODO(), lavalink.WithTrack(nextTrack)); err != nil {
		return nil, fmt.Errorf("error while skipping to next track: %w", err)
	}

	return discord.NewEmbedBuilder().
		SetTitle("Skipped Track(s)").
		SetDescription(fmt.Sprintf("Skipped %d %s.\n\nNow playing: **%s**\nby %s",
			skippedTracks,
			pluralize("track", skippedTracks),
			nextTrack.Info.Title,
			nextTrack.Info.Author)).
		SetColor(ColorSuccess).
		SetThumbnail(*nextTrack.Info.ArtworkURL), nil
}

func pluralize(word string, count int) string {
	if count == 1 {
		return word
	}
	return word + "s"
}
