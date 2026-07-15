package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	connectionString := "amqp://guest:guest@localhost:5672/"
	conn, _ := amqp.Dial(connectionString)

	ch, _ := conn.Channel()
	defer conn.Close()

	fmt.Println("Connection to Peril server successful.")

	pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, routing.GameLogSlug, routing.GameLogSlug+".*", pubsub.SimpleQueueType{Durable: true})

	gamelogic.PrintServerHelp()

outerLoop:
	for {
		userinput := gamelogic.GetInput()
		if len(userinput) == 0 {
			continue
		}
		switch userinput[0] {
		case "pause":
			fmt.Println("Sending pause message.")
			pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
		case "resume":
			fmt.Println("Sending resume message.")
			pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
		case "quit":
			fmt.Println("Exiting game.")
			break outerLoop
		default:
			fmt.Println("Command not understood.")
		}
	}

	fmt.Println("Peril shutting down.")
	conn.Close()
}
