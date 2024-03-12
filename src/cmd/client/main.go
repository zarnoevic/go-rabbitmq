package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/zarnoevic/go-rabbitmq/src/pkg/services/producerService"
)

func main() {
	fmt.Printf("Starting producer...\n")
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	rabbitMQUser := os.Getenv("RABBITMQ_USER")
	rabbitMQPassword := os.Getenv("RABBITMQ_PASSWORD")
	rabbitMQAMQPPort := os.Getenv("RABBITMQ_AMQP_PORT")
	rabbitMQAMQPHost := os.Getenv("RABBITMQ_AMQP_HOST")
	rabbitMQQueueName := os.Getenv("RABBITMQ_QUEUE_NAME")
	commandsPath := os.Getenv("COMMANDS_PATH")
	amqpURL := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitMQUser, rabbitMQPassword, rabbitMQAMQPHost, rabbitMQAMQPPort)
	fmt.Printf("Loaded .env file...\n")

	err := producerService.ProcessCSV(commandsPath, amqpURL, rabbitMQQueueName)
	if err != nil {
		fmt.Printf("Error processing CSV: %s\n", err)
	}
}
