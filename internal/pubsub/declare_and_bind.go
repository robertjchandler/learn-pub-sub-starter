package pubsub

import (
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType struct {
	Durable bool // "false" = transient
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {
	ch, _ := conn.Channel()
	PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})

	var durable, autoDelete, exclusive = false, false, false

	if queueType.Durable {
		durable, autoDelete, exclusive = true, false, false
	} else {
		durable, autoDelete, exclusive = false, true, true
	}
	queue, _ := ch.QueueDeclare(queueName, durable, autoDelete, exclusive, false, nil)

	ch.QueueBind(queueName, key, exchange, false, nil)
	return ch, queue, nil
}
