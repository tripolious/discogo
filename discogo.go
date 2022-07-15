package discogo

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"log"
	"sync"
)

func Start(ctx context.Context, wg *sync.WaitGroup, handlers []interface{}, token string) {
	wg.Add(1)
	defer wg.Done()
	log.Println("starting discord bot")

	ds, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Printf("invalid bot parameters: %v", err)
		return
	}
	defer ds.Close()

	ds.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) { log.Println("discord bot is up and running") })

	// adding handlers
	for _, handler := range handlers {
		ds.AddHandler(handler)
	}

	ds.Identify.Intents = discordgo.IntentsGuildMessages

	err = ds.Open()
	if err != nil {
		log.Printf("cant open discord bot: %s", err)
		return
	}

	select {
	case <-ctx.Done():
		log.Println("gracefully shutting down discord bot")
		err := ds.Close()
		if err != nil {
			log.Printf("error shutting down %s", err)
		}
		return
	}
}
