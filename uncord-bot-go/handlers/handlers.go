package handlers

import (
	"sync"
	"uncord-bot-go/queue"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/v3/disgolink"
)

type Handler struct {
	Client   bot.Client
	Lavalink disgolink.Client
	Queues   *queue.QueueManager
	mu       sync.Mutex
}

func NewHandler() *Handler {
	return &Handler{
		Queues: queue.NewQueueManager(),
	}
}

func (h *Handler) OnEvent(event bot.Event) {
	switch e := event.(type) {
	case *events.MessageCreate:
		h.OnMessageCreate(e)
	case *events.GuildVoiceStateUpdate:
		h.OnVoiceStateUpdate(e)
	case *events.VoiceServerUpdate:
		h.OnVoiceServerUpdate(e)
	case *events.ComponentInteractionCreate:
		h.OnComponentInteraction(e)
	case *events.ApplicationCommandInteractionCreate:
		h.HandleSlashCommand(e)
	}
}
