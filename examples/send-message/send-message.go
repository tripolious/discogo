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
	log.Printf("bot booted and is ready to send messages")

	var handlers = []interface{}{
		func(s *discordgo.Session, r *discordgo.Ready) { log.Println("discord bot is up and running") },
	}
	err = discogo.AddHandlers(handlers)
	if err != nil {
		log.Printf("unable to add handlers: %s", err)
	}

	err = discogo.SendMessage(*channel, "hello world!")
	if err != nil {
		log.Printf("unable to send message: %s", err)
	}

	select {
	case <-cancelChan:
		wg.Wait()
		return
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
