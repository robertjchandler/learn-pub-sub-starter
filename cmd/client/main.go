package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")

	connectionString := "amqp://guest:guest@localhost:5672/"
	conn, _ := amqp.Dial(connectionString)

	ch, _ := conn.Channel()
	pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})

	defer conn.Close()

	fmt.Println("Connection to Peril client successful.")

	username, _ := gamelogic.ClientWelcome()

	pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, routing.PauseKey+"."+username, routing.PauseKey, pubsub.SimpleQueueType{Durable: false})

	gamestate := gamelogic.NewGameState(username)

outerLoop:
	for {
		userinput := gamelogic.GetInput()
		if len(userinput) == 0 {
			continue
		}
		switch userinput[0] {
		case "spawn":
			gamestate.CommandSpawn(userinput)
		case "move":
			gamestate.CommandMove(userinput)
		case "status":
			gamestate.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			break outerLoop
		default:
			fmt.Println("Invalid command.")
		}
	}
}
