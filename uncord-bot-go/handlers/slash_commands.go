package handlers

import (
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
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

func (h *Handler) HandleSlashCommand(event *events.ApplicationCommandInteractionCreate) {
	switch event.Data.CommandName() {
	case "nowplaying":
		h.handleNowPlaying(event)
	case "queue":
		h.handleQueue(event)
	case "player":
		h.handlePlayer(event)
	}
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
		SetContent(fmt.Sprintf("Now playing: %s by %s", currentTrack.Info.Title, currentTrack.Info.Author)).
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
