package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"unccord-bot-go/config"
	"unccord-bot-go/handlers"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"log/slog"
)

func main() {
	slog.Info("Starting unccord-bot-go...")

	// Load configuration
	config.LoadConfig()
	if config.AppConfig.DiscordToken == "" {
		slog.Error("No Discord token provided. Check your configuration.")
		return
	}

	// Setup graceful shutdown signal handling
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	setupSignalHandler()

	// Initialize bot handlers
	b := handlers.NewHandler()

	// Create the bot client
	client, err := disgo.New(config.AppConfig.DiscordToken,
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

	// Initialize Lavalink with the loaded config
	b.Lavalink = disgolink.New(client.ApplicationID())
	if err = setupLavalink(ctx, b); err != nil {
		slog.Error("Failed to setup Lavalink", slog.Any("err", err))
		return
	}

	// Connect to Discord Gateway
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

	slog.Info("unccord-bot-go is now running. Press CTRL-C to exit.")
	<-setupSignalHandler()
}

// setupLavalink initializes Lavalink nodes based on the configuration.
func setupLavalink(ctx context.Context, b *handlers.Handler) error {
	node, err := b.Lavalink.AddNode(ctx, disgolink.NodeConfig{
		Name:     "default",
		Address:  config.AppConfig.LavalinkHost + ":" + config.AppConfig.LavalinkPort,
		Password: config.AppConfig.LavalinkPassword,
	})
	if err != nil {
		return err
	}

	version, err := node.Version(ctx)
	if err != nil {
		return err
	}

	slog.Info("Lavalink node connected", slog.String("version", version), slog.String("session_id", node.SessionID()))
	return nil
}

// setupSignalHandler handles graceful shutdown signals.
func setupSignalHandler() chan os.Signal {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	return s
}
