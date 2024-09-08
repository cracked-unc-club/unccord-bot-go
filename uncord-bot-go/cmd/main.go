package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
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

	node, err := b.Lavalink.AddNode(ctx, disgolink.NodeConfig{
		Name:     "local",
		Address:  "localhost:2333",
		Password: "youshallnotpass",
		Secure:   false,
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
