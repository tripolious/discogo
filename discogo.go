package discogo

import (
	"context"
	"errors"
	"github.com/bwmarrin/discordgo"
	"log"
	"sync"
)

var ds *discordgo.Session

func Boot(ctx context.Context, wg *sync.WaitGroup, token string) error {

	// we only allow one initialization
	if ds != nil {
		return errors.New("bot is already booted")
	}

	var err error
	ds, err = discordgo.New("Bot " + token)
	if err != nil {
		return err
	}

	ds.Identify.Intents = discordgo.IntentsGuildMessages

	err = ds.Open()
	if err != nil {
		ds.Close()
		return err
	}

	go func(wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()

		select {
		case <-ctx.Done():
			log.Println("gracefully shutting down discord bot")
			err := ds.Close()
			if err != nil {
				log.Printf("error shutting down %s", err)
			}
			return
		}
	}(wg)
	return nil
}

func AddHandlers(handlers []interface{}) error {
	// return if we don't have an active bot session
	if ds == nil {
		return errors.New("bot is not booted")
	}

	for _, handler := range handlers {
		ds.AddHandler(handler)
	}

	return nil
}

func SendMessage(channel string, message string) error {
	// return if we don't have an active bot session
	if ds == nil {
		return errors.New("bot is not booted")
	}

	_, err := ds.ChannelMessageSend(channel, message)
	if err != nil {
		return err
	}
	return nil
}
