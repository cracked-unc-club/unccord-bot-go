package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"uncord-bot-go/handlers"
	"uncord-bot-go/lavalink"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
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

	// Create the Disgo client with the appropriate intents and event listeners
	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
				gateway.IntentGuildVoiceStates,
			),
		),
	)
	if err != nil {
		slog.Error("Error while building client", slog.Any("err", err))
		return
	}

	// Initialize Lavalink client
	lavalinkClient, err := lavalink.NewClient(disgolink.NodeConfig{
		Name:     "local",
		Address:  "localhost:2333",
		Password: "youshallnotpass",
		Secure:   false,
	}, client)

	if err != nil {
		slog.Error("Error initializing Lavalink client", slog.Any("err", err))
		return
	}

	handler := handlers.NewHandler(lavalinkClient, client)
	client.AddEventListeners(handler)

	defer client.Close(context.TODO())

	if err = client.OpenGateway(context.TODO()); err != nil {
		slog.Error("Error connecting to Discord gateway", slog.Any("err", err))
		return
	}

	slog.Info("uncord-bot-go is now running. Press CTRL-C to exit.")

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
