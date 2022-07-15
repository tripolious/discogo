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

	var handlers = []interface{}{messageCreate}
	go discogo.Start(ctx, &wg, handlers, *token)

	select {
	case <-cancelChan:
		wg.Wait()
		return
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Message.ChannelID != *channel {
		return
	}

	_, err := s.ChannelMessageSend(*channel, "okokkok")
	if err != nil {
		log.Printf("cant write message %s with error %s", m.Content, err)
	}

	log.Printf("%+v", m.Content)

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
