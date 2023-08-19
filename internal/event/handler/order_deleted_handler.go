package handler

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/devfullcycle/20-CleanArch/pkg/events"
	"github.com/streadway/amqp"
)

type OrderDeletedHandler struct {
	RabbitMQChannel *amqp.Channel
}

func NewDeletedHandler(rabbitMQChannel *amqp.Channel) *OrderDeletedHandler {
	return &OrderDeletedHandler{
		RabbitMQChannel: rabbitMQChannel,
	}
}

func (h *OrderDeletedHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Order deleted: %v", event.GetPayload())
	jsonOutput, _ := json.Marshal(event.GetPayload())

	msgRabbitmq := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonOutput,
	}

	h.RabbitMQChannel.Publish(
		"deleted_exchange", // exchange
		"",                 // routing key
		false,              // mandatory
		false,              // immediate
		msgRabbitmq,        // message to publish
	)
}
