package main
import (
	"log"
	"os"
	"uncord-bot-go/config"
	"uncord-bot-go/handlers"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
)

func main() {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("No token provided. Set the DISCORD_BOT_TOKEN environment variable.")
	}

	client, err := bot.New(token, bot.WithEventListeners(&events.ListenerAdapter{
		OnMessageCreate: handlers.MessageHandler,
	}))
	if err != nil {
		log.Fatalf("Error creating bot client: %v", err)
	}
	if err = client.Open(); err != nil {
		log.Fatalf("Error connecting to Discord: %v", err)
	}
	log.println("Bot is now running...")

	// keep the bot running
	select {}
}