package graph

import "github.com/devfullcycle/20-CleanArch/internal/usecase"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	CreateOrderUseCase usecase.CreateOrderUseCase
	UpdateOrderUseCase usecase.UpdateOrderUseCase
	DeleteOrderUseCase usecase.DeleteOrderUseCase
	ListOrdersUseCase  usecase.ListOrdersUseCase
	GetOrderUseCase    usecase.GetOrderUseCase
}
