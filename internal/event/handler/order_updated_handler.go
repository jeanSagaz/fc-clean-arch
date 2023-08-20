package handler

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/jeanSagaz/fc-clean-arch/pkg/events"
	"github.com/streadway/amqp"
)

type OrderUpdatedHandler struct {
	RabbitMQChannel *amqp.Channel
}

func NewUpdatedHandler(rabbitMQChannel *amqp.Channel) *OrderUpdatedHandler {
	return &OrderUpdatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	}
}

func (h *OrderUpdatedHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Order updated: %v", event.GetPayload())
	jsonOutput, _ := json.Marshal(event.GetPayload())

	msgRabbitmq := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonOutput,
	}

	h.RabbitMQChannel.Publish(
		"updated_exchange", // exchange
		"",                 // routing key
		false,              // mandatory
		false,              // immediate
		msgRabbitmq,        // message to publish
	)
}
