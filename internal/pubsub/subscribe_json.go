package pubsub

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AckType int

const (
	SimpleQueueDurable SimpleQueueType = iota
	SimpleQueueTransient
)

const (
	Ack AckType = iota
	NackDiscard
	NackRequeue
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) AckType,
) error {
	return subscribe[T](
		conn,
		exchange,
		queueName,
		key,
		queueType,
		handler,
		func(data []byte) (T, error) {
			var target T
			err := json.Unmarshal(data, &target)
			return target, err
		},
	)
}
