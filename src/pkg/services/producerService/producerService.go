package producerService

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"

	"github.com/zarnoevic/go-rabbitmq/src/pkg/rabbitClient"
)

func ProcessCSV(filePath, amqpURL, queueName string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	client, err := rabbitClient.NewRabbitClient(amqpURL, queueName)
	if err != nil {
		return err
	}
	defer client.Close()

	var wg sync.WaitGroup
	for _, record := range records {
		wg.Add(1)
		go func(record []string) {
			defer wg.Done()
			message := record[0] + "," + record[1] + "," + record[2]
			_ = client.Publish(message)
		}(record)
	}
	wg.Wait()

	fmt.Println("All messages sent")
	return nil
}
