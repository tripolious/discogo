# discogo
discogo is a package to create a connection as a discord bot.

It's main purpose is to listen to messages and respond accordingly

## Getting Started

### Installing
```sh
go get github.com/tripolious/discogo
```

You first need to create a discord bot and a token
- [Follow this until installing your app](https://discord.com/developers/docs/getting-started#creating-an-app)


### Usage

Import the package into your project.

```go
import "github.com/tripolious/discogo
```

You can add as many handlers as you wish in your project to consume events
```go
func handler_x(s *discordgo.Session, m *discordgo.MessageCreate){}
func handler_y(s *discordgo.Session, m *discordgo.MessageCreate){}

var handlers = []interface{}{handler_x,handler_y}
go discogo.Start(ctx, &wg, handlers, *token)
```

### Example
Checkout examples for a working solution
```go
go run examples/listen-and-answer.go --token xxx --channel yyy 
```