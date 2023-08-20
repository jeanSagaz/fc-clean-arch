package usecase

import "github.com/jeanSagaz/fc-clean-arch/internal/entity"

type GetOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewGetOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
) *GetOrderUseCase {
	return &GetOrderUseCase{
		OrderRepository: OrderRepository,
	}
}

func (g *GetOrderUseCase) GetOrderById(id string) (*OrderOutputDTO, error) {
	order, err := g.OrderRepository.GetById(id)
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

func (l *ListOrdersUseCase) GetOrders() ([]OrderOutputDTO, error) {
	orders, err := l.OrderRepository.GetAll()
	if err != nil {
		return nil, err
	}

	var dto []OrderOutputDTO
	for _, o := range orders {
		dto = append(dto, OrderOutputDTO{
			ID:         o.ID,
			Price:      o.Price,
			Tax:        o.Tax,
			FinalPrice: o.FinalPrice,
		})
	}

	return dto, nil
}
