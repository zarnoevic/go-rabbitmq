package consumerService

import (
	"fmt"
	"strings"
	"sync"

	"github.com/rabbitmq/amqp091-go"
	"github.com/zarnoevic/go-rabbitmq/src/pkg/orderedmap"
)

type ConsumerService struct {
	omap *orderedmap.OrderedMap
	rwmu sync.RWMutex
}

func NewConsumerService(omap *orderedmap.OrderedMap) *ConsumerService {
	return &ConsumerService{
		omap: omap,
		rwmu: sync.RWMutex{},
	}
}

func (cs *ConsumerService) ProcessMessage(d amqp091.Delivery) {
	command := strings.Split(string(d.Body), ",")
	if len(command) < 1 {
		fmt.Println("Invalid command format")
		return
	}

	switch command[0] {
	case "add":
		cs.handleAdd(command)
	case "get":
		cs.handleGet(command)
	case "delete":
		cs.handleDelete(command)
	case "getAll":
		cs.handleGetAll()
	default:
		fmt.Printf("Unknown command: %s", command[0])
	}
}

func (cs *ConsumerService) handleAdd(command []string) {
	fmt.Printf("Add key: %s value: %s\n", command[1], command[2])
	cs.omap.Add(command[1], command[2])
}

func (cs *ConsumerService) handleGet(command []string) {
	fmt.Printf("Get key: %s\n", command[1])
	value, found := cs.omap.Get(command[1])
	if found {
		fmt.Printf("Get: %s = %s\n", command[1], value)
	} else {
		fmt.Printf("Key not found: %s\n", command[1])
	}
}

func (cs *ConsumerService) handleDelete(command []string) {
	fmt.Printf("Delete key: %s\n", command[1])
	cs.omap.Delete(command[1])
}

func (cs *ConsumerService) handleGetAll() {
	fmt.Println("Get all keys and values")
	allItems := cs.omap.GetAll()
	fmt.Printf("Retrieved all items of length: %d\n", len(allItems))
}
