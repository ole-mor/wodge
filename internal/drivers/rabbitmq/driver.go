package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"wodge/internal/services"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQDriver struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQDriver(url string) (*RabbitMQDriver, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &RabbitMQDriver{
		conn:    conn,
		channel: ch,
	}, nil
}

// Ensure RabbitMQDriver implements services.QueueService
var _ services.QueueService = (*RabbitMQDriver)(nil)

func (r *RabbitMQDriver) Publish(ctx context.Context, topic string, message []byte) error {
	// Declare queue to ensure it exists
	_, err := r.channel.QueueDeclare(
		topic, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	return r.channel.PublishWithContext(ctx,
		"",    // exchange
		topic, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
}

func (r *RabbitMQDriver) Subscribe(ctx context.Context, topic string, handler func(message []byte) error) error {
	_, err := r.channel.QueueDeclare(
		topic,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := r.channel.Consume(
		topic,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	// Run consumer in a goroutine
	// Note: This simple implementation doesn't handle graceful shutdown of the consumer properly via context yet
	go func() {
		for d := range msgs {
			if err := handler(d.Body); err != nil {
				log.Printf("Error processing message from %s: %v", topic, err)
			}
		}
	}()

	return nil
}
