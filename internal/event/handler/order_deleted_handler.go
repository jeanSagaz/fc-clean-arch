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

	setupDeleted(h.RabbitMQChannel)
	h.RabbitMQChannel.Publish(
		"deleted_exchange", // exchange
		"",                 // routing key
		false,              // mandatory
		false,              // immediate
		msgRabbitmq,        // message to publish
	)
}

func setupDeleted(rabbitMQChannel *amqp.Channel) {
	// With the instance and declare Exchanges that we can publish and subscribe to.
	err := rabbitMQChannel.ExchangeDeclare(
		"deleted_exchange", // name
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
		"deleted_queue", // queue name
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
		"deleted_queue",    // name
		"",                 // key
		"deleted_exchange", // exchange
		false,              // noWait
		nil,                // arguments - ex: amqp.Table{"alternate-exchange": "name_exchange"}
	)
}
