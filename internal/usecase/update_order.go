package usecase

import (
	"github.com/jeanSagaz/fc-clean-arch/internal/entity"
	"github.com/jeanSagaz/fc-clean-arch/pkg/events"
)

type UpdateOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	OrderUpdated    events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewUpdateOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
	OrderUpdated events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *UpdateOrderUseCase {
	return &UpdateOrderUseCase{
		OrderRepository: OrderRepository,
		OrderUpdated:    OrderUpdated,
		EventDispatcher: EventDispatcher,
	}
}

func (u *UpdateOrderUseCase) Update(input OrderRequestDTO) (OrderOutputDTO, error) {
	order := entity.Order{
		ID:    input.ID,
		Price: input.Price,
		Tax:   input.Tax,
	}
	order.CalculateFinalPrice()

	if err := u.OrderRepository.Update(&order); err != nil {
		return OrderOutputDTO{}, err
	}

	dto := OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.Price + order.Tax,
	}

	u.OrderUpdated.SetPayload(dto)
	u.EventDispatcher.Dispatch(u.OrderUpdated)

	return dto, nil
}
