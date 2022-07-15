# discogo
discogo is a package to create a connection as a discord 
bot to handle messages received in a channel and send 
messages to a specific channel.

You can use it to implement handlers for updating specific 
configurations in your app, so you can update them 
without restarting the process and also log some 
information to discord.

## Getting Started

### Installing
```sh
go get github.com/tripolious/discogo
```

You first need to create a discord bot and a token
- [Follow this until installing your app](https://discord.com/developers/docs/getting-started#creating-an-app)

### Note
The bot will run as a goroutine. To handle the shutdown gracefully you need a context and a waitgroup.  

You can find a working example (main + CreateLaunchContext func) under examples how to start it and handle the interrupt (ctr+c)  


### Usage

Import the package into your project.

```go
import "github.com/tripolious/discogo
```

Initialize the bot
```go
log.Println("starting discord bot")
err := discogo.Boot(ctx, &wg, *token)
if err != nil {
    log.Fatalf("booting failed %s", err)
}
```

You can add as many handlers as you wish in your project to consume message events
```go
// define a function you want to add
func handler_x(s *discordgo.Session, m *discordgo.MessageCreate){
	
}
// define an array of handlers
var handlers = []interface{}{
    func(s *discordgo.Session, r *discordgo.Ready) { log.Println("discord bot is up and running") }, 
    handler_x,
}

// add them
err = discogo.AddHandlers(handlers)
if err != nil {
log.Printf("unable to add handlers: %s", err)
}
```

Send a message to a discord channel
```go
err := discogo.SendMessage(*channel, "hello world!")
if err != nil {
    log.Printf("unable to send message: %s", err)
}
```

### Example
Checkout examples for a working solution
```go
go run examples/listen-and-respond/listen-and-respond.go --token xxx --channel yyy 
go run examples/send-message/send-message.go --token xxx --channel yyy 
```