package handler

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/devfullcycle/20-CleanArch/pkg/events"
	"github.com/streadway/amqp"
)

type OrderCreatedHandler struct {
	RabbitMQChannel *amqp.Channel
}

func NewOrderCreatedHandler(rabbitMQChannel *amqp.Channel) *OrderCreatedHandler {
	return &OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	}
}

func (h *OrderCreatedHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Order created: %v", event.GetPayload())
	jsonOutput, _ := json.Marshal(event.GetPayload())

	msgRabbitmq := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonOutput,
	}

	setupCreated(h.RabbitMQChannel)
	h.RabbitMQChannel.Publish(
		"created_exchange", // exchange
		"",                 // routing key
		false,              // mandatory
		false,              // immediate
		msgRabbitmq,        // message to publish
	)
}

func setupCreated(rabbitMQChannel *amqp.Channel) {
	// With the instance and declare Exchanges that we can publish and subscribe to.
	err := rabbitMQChannel.ExchangeDeclare(
		"created_exchange", // name
		"fanout",           // kind
		true,               // durable
		false,              // auto delete
		false,              // internal
		false,              // noWait
		nil,                // arguments - ex: amqp.Table{"alternate-exchange": "name_exchange"}
	)
	if err != nil {
		panic(err)
	}

	// With the instance and declare Queues that we can publish and subscribe to.
	_, err = rabbitMQChannel.QueueDeclare(
		"created_queue", // queue name
		true,            // durable
		false,           // auto delete
		false,           // exclusive
		false,           // no wait
		nil,             // arguments - ex: amqp.Table{"alternate-exchange": "name_exchange"}
	)
	if err != nil {
		panic(err)
	}

	// With the instance declare bind between Queue and Exchange.
	rabbitMQChannel.QueueBind(
		"created_queue",    // name
		"",                 // key
		"created_exchange", // exchange
		false,              // noWait
		nil,                // arguments - ex: amqp.Table{"alternate-exchange": "name_exchange"}
	)
}
