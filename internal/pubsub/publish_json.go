package pubsub

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	data, err := json.Marshal(val)
	if err != nil {
		log.Fatal(err)
	}

	pub := amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	}
	ch.PublishWithContext(context.Background(), exchange, key, false, false, pub)

	return nil
}
