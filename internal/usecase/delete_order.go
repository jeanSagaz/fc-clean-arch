package usecase

import (
	"github.com/devfullcycle/20-CleanArch/internal/entity"
	"github.com/devfullcycle/20-CleanArch/pkg/events"
)

type DeleteOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	OrderDeleted    events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewDeleteOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
	OrderDeleted events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *DeleteOrderUseCase {
	return &DeleteOrderUseCase{
		OrderRepository: OrderRepository,
		OrderDeleted:    OrderDeleted,
		EventDispatcher: EventDispatcher,
	}
}

func (d *DeleteOrderUseCase) Delete(id string) (OrderOutputDTO, error) {
	order, err := d.OrderRepository.GetById(id)
	if err != nil {
		return OrderOutputDTO{}, err
	}

	if err := d.OrderRepository.Delete(id); err != nil {
		return OrderOutputDTO{}, err
	}

	dto := OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.Price + order.Tax,
	}

	d.OrderDeleted.SetPayload(dto)
	d.EventDispatcher.Dispatch(d.OrderDeleted)

	return dto, nil
}
