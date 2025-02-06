package helpers

import (
	"encoding/json"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	RabbitMQChannel *amqp.Channel
	RabbitMQConn    *amqp.Connection
)

func InitRabbitMQ() error {
	var err error
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		err = connectRabbitMQ()
		if err == nil {
			log.Println("Successfully connected to RabbitMQ")
			return nil
		}

		log.Printf("Failed to connect to RabbitMQ (attempt %d/%d): %v\n", i+1, maxRetries, err)
		time.Sleep(5 * time.Second)
	}
	return err
}

func connectRabbitMQ() error {
	amqpUrl := getAMQPURL()

	conn, err := amqp.Dial(amqpUrl)
	if err != nil {
		return err
	}
	RabbitMQConn = conn

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	RabbitMQChannel = ch

	_, err = ch.QueueDeclare(
		"event_operations",
		true,
		false,
		false,
		false,
		nil,
	)

	return err
}

func getAMQPURL() string {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		return "amqp://guest:guest@localhost:5672/"
	}
	return url
}

func PublishToQueue(message interface{}, queueName string) error {
	q, err := RabbitMQChannel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Printf("Failed to declare queue: %v", err)
		return err
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return RabbitMQChannel.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
