package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"strconv"

	"github.com/zarnoevic/go-rabbitmq/src/pkg/orderedmap"
	"github.com/zarnoevic/go-rabbitmq/src/pkg/rabbitClient"
	"github.com/zarnoevic/go-rabbitmq/src/pkg/services/consumerService"
)

func main() {
	fmt.Println("Starting consumer...")

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	fmt.Println("Loaded .env file...")

	rabbitMQUser := os.Getenv("RABBITMQ_USER")
	rabbitMQPassword := os.Getenv("RABBITMQ_PASSWORD")
	rabbitMQAMQPPort := os.Getenv("RABBITMQ_AMQP_PORT")
	rabbitMQAMQPHost := os.Getenv("RABBITMQ_AMQP_HOST")
	rabbitMQQueueName := os.Getenv("RABBITMQ_QUEUE_NAME")

	numWorkers, err := strconv.Atoi(os.Getenv("SERVER_WORKERS"))
	if err != nil {
		log.Fatalf("Failed to parse SERVER_WORKERS: %s", err)
	}

	amqpURL := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitMQUser, rabbitMQPassword, rabbitMQAMQPHost, rabbitMQAMQPPort)

	fmt.Println("Connecting to RabbitMQ...")
	client, err := rabbitClient.NewRabbitClient(amqpURL, rabbitMQQueueName)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ client: %s", err)
	}
	defer client.Close()
	fmt.Printf("Connected to RabbitMQ with client %v\n", client)

	fmt.Println("Consuming messages...")
	messages, err := client.Consume()
	if err != nil {
		log.Fatalf("Failed to consume messages: %s", err)
	}
	fmt.Printf("Consuming messages with msgs %v\n", messages)

	omap := orderedmap.NewOrderedMap()
	fmt.Printf("Created ordered map %v\n", omap)

	processor := consumerService.NewConsumerService(omap)
	fmt.Printf("Created consumer service %v\n", processor)

	for i := 0; i < numWorkers; i++ {
		go worker(messages, processor, i)
	}

	select {}
}

func worker(messages <-chan amqp091.Delivery, processor *consumerService.ConsumerService, workerId int) {
	for d := range messages {
		//fmt.Printf("Worker %d processing message\n", workerId)
		processor.ProcessMessage(d)
	}
}
