package main

import (
	"database/sql"
	"fmt"

	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/devfullcycle/20-CleanArch/configs"

	"github.com/devfullcycle/20-CleanArch/internal/event"
	"github.com/devfullcycle/20-CleanArch/internal/event/handler"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/devfullcycle/20-CleanArch/internal/infra/database"
	"github.com/devfullcycle/20-CleanArch/internal/infra/graph"
	"github.com/devfullcycle/20-CleanArch/internal/infra/grpc/pb"
	"github.com/devfullcycle/20-CleanArch/internal/infra/grpc/service"
	"github.com/devfullcycle/20-CleanArch/internal/infra/web"
	"github.com/devfullcycle/20-CleanArch/internal/infra/web/webserver"
	"github.com/devfullcycle/20-CleanArch/pkg/events"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/streadway/amqp"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))

	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel()
	setup(rabbitMQChannel, "created", "created")
	setup(rabbitMQChannel, "updated", "updated")
	setup(rabbitMQChannel, "deleted", "deleted")
	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})
	eventDispatcher.Register("OrderUpdated", &handler.OrderUpdatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})
	eventDispatcher.Register("OrderDeleted", &handler.OrderDeletedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	updateOrderUseCase := NewUpdateOrderUseCase(db, eventDispatcher)
	deleteOrderUseCase := NewDeleteOrderUseCase(db, eventDispatcher)
	listOrdersUseCase := NewListOrdersUseCase(db)
	getOrderUseCase := NewGetOrdesUseCase(db)

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := webOrderHandlerDI(db, eventDispatcher)

	webserver.Router.Use(middleware.Logger)
	webserver.Router.Get("/order", webOrderHandler.GetAll)
	webserver.Router.Get("/order/{id}", webOrderHandler.GetById)
	webserver.Router.Post("/order", webOrderHandler.Create)
	webserver.Router.Put("/order/{id}", webOrderHandler.Update)
	webserver.Router.Delete("/order/{id}", webOrderHandler.Delete)

	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.InitOrderService(*createOrderUseCase,
		*updateOrderUseCase,
		*deleteOrderUseCase,
		*listOrdersUseCase,
		*getOrderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		UpdateOrderUseCase: *updateOrderUseCase,
		DeleteOrderUseCase: *deleteOrderUseCase,
		ListOrdersUseCase:  *listOrdersUseCase,
		GetOrderUseCase:    *getOrderUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func webOrderHandlerDI(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *web.WebOrderHandler {
	orderRepository := database.NewOrderRepository(db)
	orderCreated := event.NewOrderCreated()
	orderUpdated := event.NewOrderUpdated()
	orderDeleted := event.NewOrderDeleted()
	webOrderHandler := web.NewWebOrderHandler(eventDispatcher, orderRepository, orderCreated, orderUpdated, orderDeleted)
	return webOrderHandler
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}

func setup(rabbitMQChannel *amqp.Channel, exchange string, queue string) {
	// With the instance and declare Exchanges that we can publish and subscribe to.
	err := rabbitMQChannel.ExchangeDeclare(
		exchange+"_exchange", // name
		"fanout",             // kind
		true,                 // durable
		false,                // auto delete
		false,                // internal
		false,                // noWait
		nil,                  // arguments - ex: amqp.Table{"alternate-exchange": "name_exchange"}
	)
	if err != nil {
		panic(err)
	}

	// With the instance and declare Queues that we can publish and subscribe to.
	_, err = rabbitMQChannel.QueueDeclare(
		queue+"_queue", // queue name
		true,           // durable
		false,          // auto delete
		false,          // exclusive
		false,          // no wait
		nil,            // arguments - ex: amqp.Table{"alternate-exchange": "name_exchange"}
	)
	if err != nil {
		panic(err)
	}

	// With the instance declare bind between Queue and Exchange.
	rabbitMQChannel.QueueBind(
		queue+"_queue",       // queue name
		"",                   // key
		exchange+"_exchange", // exchange
		false,                // noWait
		nil,                  // arguments - ex: amqp.Table{"alternate-exchange": "name_exchange"}
	)
}
