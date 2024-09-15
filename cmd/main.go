package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"unccord-bot-go/config"
	"unccord-bot-go/handlers"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
	"github.com/joho/godotenv"
)

func main() {
	slog.Info("Starting unccord-bot-go...")

	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", "err", err)
	}

	config.LoadConfig()
	config.ConnectDB()

	// Debug log for Discord token (only show the first few characters)
	slog.Debug("Discord token", "token_prefix", config.AppConfig.DiscordToken[:10]+"...")

	log.Printf("Using Discord token: %s...%s", config.AppConfig.DiscordToken[:10], config.AppConfig.DiscordToken[len(config.AppConfig.DiscordToken)-10:])

	client, err := disgo.New(config.AppConfig.DiscordToken,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
				gateway.IntentGuildMessageReactions,
			),
		),
		bot.WithEventListenerFunc(handlers.OnReactionAdd),
		bot.WithEventListenerFunc(handlers.OnMessageCreate),
		bot.WithEventListenerFunc(handlers.OnReactionRemove),
	)
	if err != nil {
		log.Printf("Error creating Discord client: %v", err)
		if strings.Contains(err.Error(), "illegal base64") {
			log.Printf("Token might be malformed. Please check for any extra characters or spaces.")
		}
		os.Exit(1)
	}

	defer client.Close(context.TODO())

	if err = client.OpenGateway(context.TODO()); err != nil {
		slog.Error("Error connecting to Discord gateway", slog.Any("err", err))
		os.Exit(1)
	}

	slog.Info("unccord-bot-go is now running. Press CTRL-C to exit.")

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
