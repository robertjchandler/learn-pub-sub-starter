package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerMove(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.ArmyMove) pubsub.AckType {
	return func(move gamelogic.ArmyMove) pubsub.AckType {
		defer fmt.Print("> ")
		moveOutcome := gs.HandleMove(move)
		switch moveOutcome {
		case gamelogic.MoveOutcomeSamePlayer:
			return pubsub.Ack
		case gamelogic.MoveOutComeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			err := pubsub.PublishJSON(
				ch,
				routing.ExchangePerilTopic,
				routing.WarRecognitionsPrefix+"."+gs.GetUsername(),
				gamelogic.RecognitionOfWar{
					Attacker: move.Player,
					Defender: gs.GetPlayerSnap(),
				},
			)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		}

		fmt.Println("error: unknown move outcome")
		return pubsub.NackDiscard
	}
}

type GameLog struct {
	Exchange string
	Username string
	Message  string
}

func publishGameLog(ch *amqp.Channel, log GameLog) pubsub.AckType {
	err := pubsub.PublishGob(
		ch,
		log.Exchange,
		routing.GameLogSlug+"."+log.Username,
		log.Message,
	)
	if err != nil {
		fmt.Printf("could not publish game log: %v", err)
		return pubsub.NackRequeue
	}
	return pubsub.Ack
}

func handlerWar(gs *gamelogic.GameState, ch *amqp.Channel) func(dw gamelogic.RecognitionOfWar) pubsub.AckType {
	return func(dw gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")
		warOutcome, winner, loser := gs.HandleWar(dw)
		switch warOutcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			gl := GameLog{
				routing.ExchangePerilTopic,
				dw.Attacker.Username,
				fmt.Sprintf("%s won a war against %s", winner, loser),
			}
			return publishGameLog(ch, gl)
		case gamelogic.WarOutcomeYouWon:
			gl := GameLog{
				routing.ExchangePerilTopic,
				dw.Attacker.Username,
				fmt.Sprintf("%s won a war against %s", winner, loser),
			}
			return publishGameLog(ch, gl)
		case gamelogic.WarOutcomeDraw:
			gl := GameLog{
				routing.ExchangePerilTopic,
				dw.Attacker.Username,
				fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser),
			}
			return publishGameLog(ch, gl)
		}

		fmt.Println("error: unknown war outcome")
		return pubsub.NackDiscard
	}
}

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.AckType {
	return func(ps routing.PlayingState) pubsub.AckType {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
		return pubsub.Ack
	}
}
