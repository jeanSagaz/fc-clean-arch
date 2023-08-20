package usecase

import "github.com/jeanSagaz/fc-clean-arch/internal/entity"

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrdersUseCase(
	OrderRepository entity.OrderRepositoryInterface,
) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: OrderRepository,
	}
}

func (l *ListOrdersUseCase) GetOrderById(id string) (*OrderOutputDTO, error) {
	order, err := l.OrderRepository.GetById(id)
	if err != nil {
		return nil, err
	}

	dto := &OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}

	return dto, nil
}
