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

    queue := h.Queues.Get(guildID)
    player := h.Lavalink.Player(guildID)

    var loadError error
    var trackLoaded bool
    var addedTracks []lavalink.Track

    isPlaying := player != nil && player.Track() != nil

    startIfNeeded := func(track lavalink.Track) {
        if !isPlaying {
            err := h.playTrack(guildID, track)
            if err != nil {
                slog.Error("Failed to play track", slog.Any("err", err))
            } else {
                trackLoaded = true
                isPlaying = true
                // Create the player control panel after the track starts playing
                go h.createControlPanel(commandChannelID, guildID)
            }
        } else {
            queue.Add(track)
            slog.Info("Added track to queue", "track", track.Info.Title, "position", len(queue.Tracks), "guildID", guildID)
            trackLoaded = true
        }
    }

    h.Lavalink.BestNode().LoadTracksHandler(context.TODO(), url, disgolink.NewResultHandler(
        func(track lavalink.Track) {
            slog.Info("Single track loaded", "title", track.Info.Title, "guildID", guildID)
            addedTracks = append(addedTracks, track)
            startIfNeeded(track)
        },
        func(playlist lavalink.Playlist) {
            slog.Info("Playlist loaded", "trackCount", len(playlist.Tracks), "guildID", guildID)
            addedTracks = append(addedTracks, playlist.Tracks...)
            for _, track := range playlist.Tracks {
                queue.Add(track)
            }
            if len(playlist.Tracks) > 0 {
                startIfNeeded(playlist.Tracks[0])
            }
        },
        func(tracks []lavalink.Track) {
            slog.Info("Search results loaded", "trackCount", len(tracks), "guildID", guildID)
            if len(tracks) > 0 {
                addedTracks = append(addedTracks, tracks...)
                for _, track := range tracks {
                    queue.Add(track)
                }
                startIfNeeded(tracks[0])
            }
            trackLoaded = true
        },
        func() {
            slog.Error("No matches found", "url", url, "guildID", guildID)
            loadError = fmt.Errorf("no matches found for URL: %s", url)
        },
        func(err error) {
            slog.Error("Error loading track", slog.Any("err", err), "url", url, "guildID", guildID)
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
    var embed *discord.EmbedBuilder
    if len(addedTracks) == 1 {
        track := addedTracks[0]
        if !isPlaying {
            embed = discord.NewEmbedBuilder().
                SetTitle("Now Playing").
                SetDescription(fmt.Sprintf("**%s**", track.Info.Title)).
                SetColor(ColorSuccess).
                SetThumbnail(*track.Info.ArtworkURL)
        } else {
            embed = discord.NewEmbedBuilder().
                SetTitle("Added to Queue").
                SetDescription(fmt.Sprintf("**%s**\n\nPosition in queue: %d",
                    track.Info.Title,
                    len(queue.Tracks))).
                SetColor(ColorInfo).
                SetThumbnail(*track.Info.ArtworkURL)
        }
    } else {
        embed = discord.NewEmbedBuilder().
            SetTitle("Playlist Added to Queue").
            SetDescription(fmt.Sprintf("Added %d tracks to the queue", len(addedTracks))).
            SetColor(ColorInfo)
    }

    _, err = h.Client.Rest().CreateMessage(commandChannelID, discord.NewMessageCreateBuilder().
        SetEmbeds(embed.Build()).
        Build())
    if err != nil {
        slog.Error("Failed to send queue message", slog.Any("err", err))
    }

    return nil
}

func (h *Handler) createControlPanel(channelID, guildID snowflake.ID) {
    player := h.Lavalink.ExistingPlayer(guildID)
    if player == nil {
        slog.Error("No active player found", "guildID", guildID)
        return
    }

    currentTrack := player.Track()
    if currentTrack == nil {
        slog.Error("No current track found", "guildID", guildID)
        return
    }

    queue := h.Queues.Get(guildID)
    queueInfo := fmt.Sprintf("Next in queue: %d", len(queue.Tracks))

    _, err := h.Client.Rest().CreateMessage(channelID, discord.NewMessageCreateBuilder().
        SetContent("").
        SetEmbeds(discord.NewEmbedBuilder().
            SetTitle("Now Playing").
            SetDescription(fmt.Sprintf("**%s**\nby *%s*\n\n%s", currentTrack.Info.Title, currentTrack.Info.Author, queueInfo)).
            SetColor(ColorInfo).
            SetThumbnail(*currentTrack.Info.ArtworkURL).
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
    } else {
        slog.Info("Control panel created successfully", "guildID", guildID)
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
    queue.Tracks = queue.Tracks[1:] // Remove the first track from the queue

    player := h.Lavalink.Player(guildID)
    err := player.Update(context.TODO(), lavalink.WithTrack(nextTrack), lavalink.WithPaused(false))
    if err != nil {
        slog.Error("Failed to play next track", slog.Any("err", err))
        // If we fail to play this track, try the next one
        h.playNextTrack(guildID)
    } else {
        slog.Info("Now playing next track", "title", nextTrack.Info.Title, "guildID", guildID)
        // Re-add the track end listener for the new track
        h.addTrackEndListener(player, guildID)
        
        // Ensure the player state is updated
        updatedPlayer := h.Lavalink.ExistingPlayer(guildID)
        if updatedPlayer != nil {
            currentTrack := updatedPlayer.Track()
            if currentTrack != nil {
                slog.Info("Player state updated", "currentTrack", currentTrack.Info.Title, "guildID", guildID)
            } else {
                slog.Error("Player state not updated correctly", "guildID", guildID)
            }
        }
    }
}

func (h *Handler) addTrackEndListener(player disgolink.Player, guildID snowflake.ID) {
    player.Node().Lavalink().AddListeners(disgolink.NewListenerFunc(func(p disgolink.Player, event lavalink.TrackEndEvent) {
        if event.Reason == lavalink.TrackEndReasonFinished {
            slog.Info("Track ended, playing next track", "guildID", guildID)
            go h.playNextTrack(guildID)  // Use a goroutine to avoid blocking
        }
    }))
}

func (h *Handler) playTrack(guildID snowflake.ID, track lavalink.Track) error {
    player := h.Lavalink.Player(guildID)
    err := player.Update(context.TODO(), lavalink.WithTrack(track), lavalink.WithPaused(false))
    if err != nil {
        slog.Error("Error updating player", slog.Any("err", err))
        return err
    }

    h.addTrackEndListener(player, guildID)

    return nil
}

func (h *Handler) handleSkipButton(event *events.ComponentInteractionCreate) {
    player := h.Lavalink.ExistingPlayer(*event.GuildID())
    if player == nil {
        event.CreateMessage(discord.NewMessageCreateBuilder().
            SetContent("No active player found.").
            SetEphemeral(true).
            Build())
        return
    }

    h.playNextTrack(*event.GuildID())

    currentTrack := player.Track()
    if currentTrack == nil {
        event.CreateMessage(discord.NewMessageCreateBuilder().
            SetContent("Skipped the last track. The queue is now empty.").
            SetEphemeral(true).
            Build())
        return
    }

    event.CreateMessage(discord.NewMessageCreateBuilder().
        SetContent(fmt.Sprintf("Skipped to next track: **%s**", currentTrack.Info.Title)).
        SetEphemeral(true).
        Build())
}

func (h *Handler) skipTracks(guildID snowflake.ID, amount int) (*discord.EmbedBuilder, error) {
    player := h.Lavalink.ExistingPlayer(guildID)
    queue := h.Queues.Get(guildID)
    if player == nil {
        return nil, fmt.Errorf("no player found")
    }

    skippedTracks := min(amount, len(queue.Tracks))
    queue.Tracks = queue.Tracks[skippedTracks:]

    if len(queue.Tracks) == 0 {
        // If we've skipped all tracks, stop the player
        if err := player.Update(context.TODO(), lavalink.WithNullTrack()); err != nil {
            return nil, fmt.Errorf("error while stopping track: %w", err)
        }
        // Reset the player state
        h.Lavalink.RemovePlayer(guildID)
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
        SetDescription(fmt.Sprintf("Skipped %d %s.\n\nNow playing: **%s**",
            skippedTracks,
            pluralize("track", skippedTracks),
            nextTrack.Info.Title)).
        SetColor(ColorSuccess).
        SetThumbnail(*nextTrack.Info.ArtworkURL), nil
}

func pluralize(word string, count int) string {
    if count == 1 {
        return word
    }
    return word + "s"
}