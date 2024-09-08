package handlers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

var commands = []discord.ApplicationCommandCreate{
	discord.SlashCommandCreate{
		Name:        "nowplaying",
		Description: "Show the currently playing song",
	},
	discord.SlashCommandCreate{
		Name:        "queue",
		Description: "Show the current music queue",
	},
	discord.SlashCommandCreate{
		Name:        "player",
		Description: "Control the music player",
	},
	discord.SlashCommandCreate{
		Name:        "skip",
		Description: "Skip the currently playing song",
	},
	discord.SlashCommandCreate{
		Name:        "pause",
		Description: "Pause the music player",
	},
	discord.SlashCommandCreate{
		Name:        "resume",
		Description: "Resume the music player",
	},
	discord.SlashCommandCreate{
		Name:        "leave",
		Description: "Leave the voice channel",
	},
}

func (h *Handler) HandleSlashCommand(event *events.ApplicationCommandInteractionCreate) {
	switch event.Data.CommandName() {
	case "nowplaying":
		h.handleNowPlaying(event)
	case "queue":
		h.handleQueue(event)
	case "player":
		h.handlePlayer(event)
	case "skip":
		h.handleSkip(event)
	case "pause":
		h.handlePause(event)
	case "resume":
		h.handleResume(event)
	case "leave":
		h.handleLeave(event)
	}
}

func (h *Handler) RegisterCommands(client bot.Client) error {
	_, err := client.Rest().SetGlobalCommands(client.ApplicationID(), commands)
	return err
}

func (h *Handler) RegisterGuildCommands(client bot.Client, guildID snowflake.ID) error {
	_, err := client.Rest().SetGuildCommands(client.ApplicationID(), guildID, commands)
	if err != nil {
		slog.Error("Failed to register guild commands", slog.Any("err", err), slog.Any("guildID", guildID))
		return err
	}

	slog.Info("Successfully registered guild commands", slog.Any("guildID", guildID))
	return nil
}

func (h *Handler) handleNowPlaying(event *events.ApplicationCommandInteractionCreate) {
	player := h.Lavalink.Player(*event.GuildID())
	queue := h.Queues.Get(*event.GuildID())
	if player == nil || len(queue.Tracks) == 0 {
		event.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("No song is currently playing.").
			SetEphemeral(true).
			Build())
		return
	}

	currentTrack := queue.Tracks[0]
	event.CreateMessage(discord.NewMessageCreateBuilder().
		SetContent(fmt.Sprintf("Now playing: **%s**", currentTrack.Info.Title)).
		SetEphemeral(true).
		Build())
}

func (h *Handler) handleQueue(event *events.ApplicationCommandInteractionCreate) {
	queue := h.Queues.Get(*event.GuildID())
	if len(queue.Tracks) == 0 {
		event.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("The queue is currently empty.").
			SetEphemeral(true).
			Build())
		return
	}

	var queueList string
	for i, track := range queue.Tracks {
		queueList += fmt.Sprintf("%d. %s by %s\n", i+1, track.Info.Title, track.Info.Author)
	}

	event.CreateMessage(discord.NewMessageCreateBuilder().
		SetContent("Current queue:").
		SetEmbeds(discord.NewEmbedBuilder().
			SetTitle("Music Queue").
			SetDescription(queueList).
			Build()).
		SetEphemeral(true).
		Build())
}

func (h *Handler) handlePlayer(event *events.ApplicationCommandInteractionCreate) {
	h.createControlPanel(event.Channel().ID(), *event.GuildID())
}

func (h *Handler) handleSkip(event *events.ApplicationCommandInteractionCreate) {
	amount := 1
	if data, ok := event.SlashCommandInteractionData().OptInt("amount"); ok {
		amount = data
	}

	embed, err := h.skipTracks(*event.GuildID(), amount)
	if err != nil {
		event.CreateMessage(discord.NewMessageCreateBuilder().
			SetEmbeds(discord.NewEmbedBuilder().
				SetDescription(fmt.Sprintf("Error: %s", err)).
				SetColor(ColorError).
				Build()).
			Build())
		return
	}

	event.CreateMessage(discord.NewMessageCreateBuilder().
		SetEmbeds(embed.Build()).
		Build())
}

func (h *Handler) handlePause(event *events.ApplicationCommandInteractionCreate) {
	player := h.Lavalink.Player(*event.GuildID())
	if player == nil {
		event.CreateMessage(discord.NewMessageCreateBuilder().
			SetEmbeds(discord.NewEmbedBuilder().
				SetDescription("No music player found for this guild.").
				SetColor(ColorWarning).
				Build()).
			Build())
		return
	}

	if err := player.Update(context.TODO(), lavalink.WithPaused(true)); err != nil {
		event.CreateMessage(discord.NewMessageCreateBuilder().
			SetEmbeds(discord.NewEmbedBuilder().
				SetDescription(fmt.Sprintf("Error: %s", err)).
				SetColor(ColorError).
				Build()).
			Build())
		return
	}

	event.CreateMessage(discord.NewMessageCreateBuilder().
		SetEmbeds(discord.NewEmbedBuilder().
			SetDescription("Music paused.").
			SetColor(ColorSuccess).
			Build()).
		Build())
}

func (h *Handler) handleResume(event *events.ApplicationCommandInteractionCreate) {
	player := h.Lavalink.Player(*event.GuildID())
	if player == nil {
		event.CreateMessage(discord.NewMessageCreateBuilder().
			SetEmbeds(discord.NewEmbedBuilder().
				SetDescription("No music player found for this guild.").
				SetColor(ColorWarning).
				Build()).
			Build())
		return
	}

	if err := player.Update(context.TODO(), lavalink.WithPaused(false)); err != nil {
		event.CreateMessage(discord.NewMessageCreateBuilder().
			SetEmbeds(discord.NewEmbedBuilder().
				SetDescription(fmt.Sprintf("Error: %s", err)).
				SetColor(ColorError).
				Build()).
			Build())
		return
	}

	event.CreateMessage(discord.NewMessageCreateBuilder().
		SetEmbeds(discord.NewEmbedBuilder().
			SetDescription("Music resumed.").
			SetColor(ColorSuccess).
			Build()).
		Build())
}

// func (b *Bot) disconnect(event *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) error {
// 	player := b.Lavalink.ExistingPlayer(*event.GuildID())
// 	if player == nil {
// 		return event.CreateMessage(discord.MessageCreate{
// 			Content: "No player found",
// 		})
// 	}

// 	if err := b.Client.UpdateVoiceState(context.TODO(), *event.GuildID(), nil, false, false); err != nil {
// 		return event.CreateMessage(discord.MessageCreate{
// 			Content: fmt.Sprintf("Error while disconnecting: `%s`", err),
// 		})
// 	}

// 	return event.CreateMessage(discord.MessageCreate{
// 		Content: "Player disconnected",
// 	})
// }

func (h *Handler) handleLeave(event *events.ApplicationCommandInteractionCreate) {
	player := h.Lavalink.Player(*event.GuildID())
	if player == nil {
		event.CreateMessage(discord.NewMessageCreateBuilder().
			SetEmbeds(discord.NewEmbedBuilder().
				SetDescription("No music player found for this guild.").
				SetColor(ColorWarning).
				Build()).
			Build())
		return
	}

	if err := player.Destroy(context.TODO()); err != nil {
		event.CreateMessage(discord.NewMessageCreateBuilder().
			SetEmbeds(discord.NewEmbedBuilder().
				SetDescription(fmt.Sprintf("Error: %s", err)).
				SetColor(ColorError).
				Build()).
			Build())
		return
	}

	// Clear the queue
	h.Queues.Get(*event.GuildID()).Tracks = nil

	// Disconnect bot from voice channel
	if err := h.Client.UpdateVoiceState(context.TODO(), *event.GuildID(), nil, false, false); err != nil {
		event.CreateMessage(discord.NewMessageCreateBuilder().
			SetEmbeds(discord.NewEmbedBuilder().
				SetDescription(fmt.Sprintf("Error while disconnecting: `%s`", err)).
				SetColor(ColorError).
				Build()).
			Build())
		return
	}

	event.CreateMessage(discord.NewMessageCreateBuilder().
		SetEmbeds(discord.NewEmbedBuilder().
			SetDescription("Left the voice channel and cleared the queue.").
			SetColor(ColorSuccess).
			Build()).
		Build())
}
