//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/google/wire"
	"github.com/jeanSagaz/fc-clean-arch/internal/entity"
	"github.com/jeanSagaz/fc-clean-arch/internal/event"
	"github.com/jeanSagaz/fc-clean-arch/internal/infra/database"
	"github.com/jeanSagaz/fc-clean-arch/internal/usecase"
	"github.com/jeanSagaz/fc-clean-arch/pkg/events"
)

var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

var setEventDispatcherDependency = wire.NewSet(
	events.NewEventDispatcher,
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
	wire.Bind(new(events.EventDispatcherInterface), new(*events.EventDispatcher)),
)

var setOrderCreatedEvent = wire.NewSet(
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
)

var setOrderUpdatedEvent = wire.NewSet(
	event.NewOrderUpdated,
	wire.Bind(new(events.EventInterface), new(*event.OrderUpdated)),
)

var setOrderDeletedEvent = wire.NewSet(
	event.NewOrderDeleted,
	wire.Bind(new(events.EventInterface), new(*event.OrderDeleted)),
)

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		usecase.NewCreateOrderUseCase,
	)
	return &usecase.CreateOrderUseCase{}
}

func NewUpdateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.UpdateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderUpdatedEvent,
		usecase.NewUpdateOrderUseCase,
	)
	return &usecase.UpdateOrderUseCase{}
}

func NewDeleteOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.DeleteOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderDeletedEvent,
		usecase.NewDeleteOrderUseCase,
	)
	return &usecase.DeleteOrderUseCase{}
}

func NewListOrdersUseCase(db *sql.DB) *usecase.ListOrdersUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		usecase.NewListOrdersUseCase,
	)
	return &usecase.ListOrdersUseCase{}
}

func NewGetOrdesUseCase(db *sql.DB) *usecase.GetOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		usecase.NewGetOrderUseCase,
	)
	return &usecase.GetOrderUseCase{}
}

// func NewWebOrderHandler(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *web.WebOrderHandler {
// 	wire.Build(
// 		setOrderRepositoryDependency,
// 		setOrderCreatedEvent,
//      setOrderUpdatedEvent
//      setOrderDeletedEvent
// 		web.NewWebOrderHandler,
// 	)
// 	return &web.WebOrderHandler{}
// }
