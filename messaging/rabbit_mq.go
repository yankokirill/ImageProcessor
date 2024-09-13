package messaging

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	. "hw/models"
	"log"
)

var _ Producer = ProducerRMQ{}
var _ Consumer = ConsumerRMQ{}

type Producer interface {
	Publish(task *Task) error
}

type Consumer interface {
	Consume() <-chan amqp.Delivery
}

type ProducerRMQ struct {
	ch *amqp.Channel
}

type ConsumerRMQ struct {
	ch *amqp.Channel
}

func NewProducerRMQ(rabbitMQAddr string) ProducerRMQ {
	return ProducerRMQ{createChannel(rabbitMQAddr)}
}

func NewConsumerRMQ(rabbitMQAddr string) ConsumerRMQ {
	return ConsumerRMQ{createChannel(rabbitMQAddr)}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func createChannel(rabbitMQAddr string) *amqp.Channel {
	conn, err := amqp.Dial(rabbitMQAddr)
	failOnError(err, "failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "failed to open a channel")

	_, err = ch.QueueDeclare(
		"task_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "failed to declare a queue")
	return ch
}

func (c ConsumerRMQ) Consume() <-chan amqp.Delivery {
	msgs, err := c.ch.Consume(
		"task_queue",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "failed to register a consumer")
	return msgs
}

func (b ProducerRMQ) Publish(task *Task) error {
	body, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	err = b.ch.Publish(
		"",
		"task_queue",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish task: %w", err)
	}

	return nil
}
