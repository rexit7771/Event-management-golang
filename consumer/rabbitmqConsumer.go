package consumer

import (
	"encoding/json"
	"event-management/helpers"
	"event-management/structs"
	"log"
	"time"
)

func ConsumeEventOperations() {
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		if helpers.RabbitMQChannel == nil {
			log.Printf("RabbitMQ channel not initialized, retrying (%d/%d)...\n", i+1, maxRetries)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	if helpers.RabbitMQChannel == nil {
		log.Printf("Failed to initialize RabbitMQ channel after %d attempts\n", maxRetries)
		return
	}

	ch := helpers.RabbitMQChannel
	q, err := ch.QueueDeclare(
		"event_operations",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var message structs.EventMessage
			json.Unmarshal(d.Body, &message)

			switch message.Action {
			case "create":
				log.Printf("Event created: %v", message.EventID)
			case "update":
				log.Printf("Event updated: %v", message.EventID)
			case "approval":
				log.Printf("Event approval updated: %v", message.EventID)
			case "delete":
				log.Printf("Event deleted: %v", message.EventID)
			}
		}
	}()
	log.Printf("Consumer started. Waiting for messages...")
	<-forever
}
