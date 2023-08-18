package web

import (
	"encoding/json"
	"net/http"

	"github.com/devfullcycle/20-CleanArch/internal/entity"
	"github.com/devfullcycle/20-CleanArch/internal/usecase"
	"github.com/devfullcycle/20-CleanArch/pkg/events"
	"github.com/go-chi/chi/v5"
)

type WebOrderHandler struct {
	EventDispatcher   events.EventDispatcherInterface
	OrderRepository   entity.OrderRepositoryInterface
	OrderCreatedEvent events.EventInterface
	OrderUpdatedEvent events.EventInterface
	OrderDeletedEvent events.EventInterface
}

func NewWebOrderHandler(
	EventDispatcher events.EventDispatcherInterface,
	OrderRepository entity.OrderRepositoryInterface,
	OrderCreatedEvent events.EventInterface,
	OrderUpdatedEvent events.EventInterface,
	OrderDeletedEvent events.EventInterface,
) *WebOrderHandler {
	return &WebOrderHandler{
		EventDispatcher:   EventDispatcher,
		OrderRepository:   OrderRepository,
		OrderCreatedEvent: OrderCreatedEvent,
		OrderUpdatedEvent: OrderUpdatedEvent,
		OrderDeletedEvent: OrderDeletedEvent,
	}
}

func (h *WebOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecase.OrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createOrder := usecase.NewCreateOrderUseCase(h.OrderRepository, h.OrderCreatedEvent, h.EventDispatcher)
	output, err := createOrder.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}

func (h *WebOrderHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	us := usecase.NewListOrdersUseCase(h.OrderRepository)
	order, err := us.GetOrderById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var dto usecase.OrderRequestDTO
	err = json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if order.ID != dto.ID {
		http.Error(w, "Os id do banco de dados não corresponde com o id do request", http.StatusBadRequest)
		return
	}

	updateOrder := usecase.NewUpdateOrderUseCase(h.OrderRepository, h.OrderUpdatedEvent, h.EventDispatcher)
	output, err := updateOrder.Update(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}

func (h *WebOrderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	us := usecase.NewListOrdersUseCase(h.OrderRepository)
	_, err := us.GetOrderById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	deleteOrder := usecase.NewDeleteOrderUseCase(h.OrderRepository, h.OrderDeletedEvent, h.EventDispatcher)
	_, err = deleteOrder.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *WebOrderHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	us := usecase.NewGetOrderUseCase(h.OrderRepository)
	order, err := us.GetOrderById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}

func (h *WebOrderHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	us := usecase.NewListOrdersUseCase(h.OrderRepository)
	output, err := us.GetOrders()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
