package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"  // Import godotenv
	"github.com/streadway/amqp" // Import amqp
)

// Define a custom error type for our RabbitMQ client
type RMQError struct {
	Err  error
	Desc string
}

func (e *RMQError) Error() string {
	return e.Desc + ": " + e.Err.Error()
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Get RabbitMQ connection URL from environment variable
	amqpURL := os.Getenv("RABBITMQ_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@localhost:5672/" // Default if not set
		log.Printf("RABBITMQ_URL not set, using default: %s", amqpURL)
	}

	exchangeName := os.Getenv("EXCHANGE_NAME")
	if exchangeName == "" {
		exchangeName = "logs_exchange" // Default exchange name
		log.Printf("EXCHANGE_NAME not set, using default: %s", exchangeName)
	}

	queueName := os.Getenv("QUEUE_NAME")
	if queueName == "" {
		queueName = "logs_queue" // Default queue name
		log.Printf("QUEUE_NAME not set, using default: %s", queueName)
	}

	// Configurable retry count and delay
	retryCount := 5
	retryDelay := 5 * time.Second

	// Establish connection
	for i := 0; i < retryCount; i++ {
		conn, err := amqp.Dial(amqpURL)
		if err != nil {
			log.Printf("Failed to connect to RabbitMQ (attempt %d/%d): %v", i+1, retryCount, err)
			if i < retryCount-1 {
				log.Println("Retrying in 5 seconds...")
				time.Sleep(retryDelay)
			} else {
				log.Fatalf("Could not connect to RabbitMQ after %d attempts: %v", retryCount, err)
			}
			continue
		}
		defer conn.Close()

		// Open a channel
		ch, err := conn.Channel()
		if err != nil {
			log.Printf("Failed to open a channel: %v", err)
			continue
		}
		defer ch.Close()

		// Declare a durable, direct exchange
		err = ch.ExchangeDeclare(
			exchangeName, // name
			"direct",     // type
			true,         // durable
			false,        // auto-deleted
			false,        // internal
			false,        // no-wait
			nil,          // arguments
		)
		if err != nil {
			log.Printf("Failed to declare exchange: %v", err)
			continue
		}

		// Declare a durable queue
		q, err := ch.QueueDeclare(
			queueName, // name
			true,      // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		if err != nil {
			log.Printf("Failed to declare queue: %v", err)
			continue
		}

		// Bind the queue to the exchange
		err = ch.QueueBind(
			q.Name,       // queue name
			"",           // routing key
			exchangeName, // exchange
			false,        // no-wait
			nil,          // args
		)
		if err != nil {
			log.Printf("Failed to bind queue: %v", err)
			continue
		}

		// Message body (You might get this from user input or another source)
		body := "This is a log message"

		// Enable publisher confirms
		if err := ch.Confirm(false); err != nil {
			log.Fatalf("Failed to enable publisher confirms: %v", err)
		}
		confirmChan := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

		// Publish a message
		err = ch.Publish(
			exchangeName, // exchange
			"",           // routing key
			false,        // mandatory
			false,        // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			},
		)
		if err != nil {
			log.Printf("Failed to publish a message: %v", err)
			continue
		}

		// Wait for confirmation
		select {
		case confirm := <-confirmChan:
			if !confirm.Ack {
				log.Printf("Message not acknowledged by RabbitMQ: %v", confirm.DeliveryTag)
			} else {
				log.Printf("Message published and acknowledged: %s", body)
			}
		case <-time.After(5 * time.Second):
			log.Println("No confirmation received within the timeout period")
		}

		// Exit the loop after successful operations
		break
	}
}
