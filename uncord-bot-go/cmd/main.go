package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"uncord-bot-go/handlers"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgolink/v3/disgolink"
)

func main() {
	slog.Info("Starting uncord-bot-go...")

	token := os.Getenv("DISCORD_TOKEN")
	lavalinkHost := os.Getenv("LAVALINK_HOST")
	lavalinkPort := os.Getenv("LAVALINK_PORT")
	lavalinkPassword := os.Getenv("LAVALINK_PASSWORD")
	lavalinkSecure, _ := strconv.ParseBool(os.Getenv("LAVALINK_SECURE"))

	if token == "" {
		slog.Error("No token provided. Set the DISCORD_TOKEN environment variable.")
		return
	}

	b := handlers.NewHandler()

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentGuilds, gateway.IntentGuildVoiceStates, gateway.IntentGuildMessages, gateway.IntentMessageContent),
		),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagVoiceStates),
		),
		bot.WithEventListeners(b),
	)
	if err != nil {
		slog.Error("Error while building client", slog.Any("err", err))
		return
	}
	b.Client = client

	b.Lavalink = disgolink.New(client.ApplicationID())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = client.OpenGateway(ctx); err != nil {
		slog.Error("Error connecting to Discord gateway", slog.Any("err", err))
		return
	}
	defer client.Close(context.TODO())

	// Register commands after connecting to the gateway
	if err = b.RegisterCommands(client); err != nil {
		slog.Error("Failed to register commands", slog.Any("err", err))
		return
	}

	// Register guild commands after connecting to the gateway (useful for testing in a specific guild without re-registering)
	// b.RegisterGuildCommands(client, snowflake.ID(YOUR_GUILD_ID))

	node, err := b.Lavalink.AddNode(ctx, disgolink.NodeConfig{
		Name:     "default",
		Address:  lavalinkHost + ":" + lavalinkPort,
		Password: lavalinkPassword,
		Secure:   lavalinkSecure,
	})
	if err != nil {
		slog.Error("Failed to add node", slog.Any("err", err))
		return
	}

	version, err := node.Version(ctx)
	if err != nil {
		slog.Error("Failed to get node version", slog.Any("err", err))
		return
	}

	slog.Info("uncord-bot-go is now running. Press CTRL-C to exit.", slog.String("node_version", version), slog.String("node_session_id", node.SessionID()))

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
