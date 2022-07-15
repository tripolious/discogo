package main

import (
	"context"
	"flag"
	"github.com/bwmarrin/discordgo"
	"github.com/tripolious/discogo"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var wg sync.WaitGroup

// params
var (
	token   = flag.String("token", "", "bot token")
	channel = flag.String("channel", "", "channel bot listens to new messages")
)

func main() {
	flag.Parse()
	ctx, cancelFunc, cancelChan := CreateLaunchContext()
	defer cancelFunc()

	log.Println("starting discord bot")
	err := discogo.Boot(ctx, &wg, *token)
	if err != nil {
		log.Fatalf("booting failed %s", err)
	}

	log.Println("add handlers so the bot can handle messages")
	var handlers = []interface{}{
		messageCreate,
		func(s *discordgo.Session, r *discordgo.Ready) { log.Println("discord bot is up and running") },
	}
	err = discogo.AddHandlers(handlers)
	if err != nil {
		log.Printf("unable to add handlers: %s", err)
	}

	log.Printf("bot booted and can consume handlers and send messages")

	select {
	case <-cancelChan:
		wg.Wait()
		return
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// we skip messages the bot did
	if m.Author.ID == s.State.User.ID {
		return
	}

	// we only want to consume messages from the defined channel
	if m.Message.ChannelID != *channel {
		return
	}

	// example of how we can just write stuff to the channel
	_, err := s.ChannelMessageSend(*channel, "okokkok")
	if err != nil {
		log.Printf("cant write message %s with error %s", m.Content, err)
	}

	// this is the message we consumed
	log.Printf("%+v", m.Content)

	// example of how we can respond to the message
	_, err = s.ChannelMessageSendReply(*channel, "yo", m.Reference())
	if err != nil {
		log.Printf("cant reply to Message %s with error %s", m.Content, err)
	}
}

func CreateLaunchContext() (context.Context, func(), chan bool) {
	interruptChan := make(chan os.Signal, 1)
	canceledChanChan := make(chan bool, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		defer close(interruptChan)
		<-interruptChan
		cancelCtx()
		canceledChanChan <- true
	}()
	cancel := func() {
		cancelCtx()
		close(canceledChanChan)
	}
	return ctx, cancel, canceledChanChan
}
