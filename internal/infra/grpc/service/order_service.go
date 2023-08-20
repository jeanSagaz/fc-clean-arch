package service

import (
	"context"

	"github.com/jeanSagaz/fc-clean-arch/internal/infra/grpc/pb"
	"github.com/jeanSagaz/fc-clean-arch/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	UpdateOrderUseCase usecase.UpdateOrderUseCase
	DeleteOrderUseCase usecase.DeleteOrderUseCase
	ListOrdersUseCase  usecase.ListOrdersUseCase
	GetOrderUseCase    usecase.GetOrderUseCase
}

func InitOrderService(createOrderUseCase usecase.CreateOrderUseCase,
	updateOrderUseCase usecase.UpdateOrderUseCase,
	deleteOrderUseCase usecase.DeleteOrderUseCase,
	listOrdersUseCase usecase.ListOrdersUseCase,
	getOrderUseCase usecase.GetOrderUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		UpdateOrderUseCase: updateOrderUseCase,
		DeleteOrderUseCase: deleteOrderUseCase,
		ListOrdersUseCase:  listOrdersUseCase,
		GetOrderUseCase:    getOrderUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.OrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) DeleteOrder(ctx context.Context, in *pb.OrderByIdRequest) (*pb.OrderResponse, error) {
	_, err := s.GetOrderUseCase.GetOrderById(in.Id)
	if err != nil {
		return nil, err
	}

	dto, err := s.DeleteOrderUseCase.Delete(in.Id)
	if err != nil {
		return nil, err
	}

	return &pb.OrderResponse{
		Id:         dto.ID,
		Price:      float32(dto.Price),
		Tax:        float32(dto.Tax),
		FinalPrice: float32(dto.FinalPrice),
	}, nil
}

func (s *OrderService) GetOrderById(ctx context.Context, in *pb.OrderByIdRequest) (*pb.OrderResponse, error) {
	dto, err := s.GetOrderUseCase.GetOrderById(in.Id)
	if err != nil {
		return nil, err
	}

	return &pb.OrderResponse{
		Id:         dto.ID,
		Price:      float32(dto.Price),
		Tax:        float32(dto.Tax),
		FinalPrice: float32(dto.FinalPrice),
	}, nil
}

func (s *OrderService) ListOrders(context.Context, *pb.BlankRequest) (*pb.OrderListResponse, error) {
	orders, err := s.ListOrdersUseCase.GetOrders()
	if err != nil {
		return nil, err
	}

	var ordersResponse []*pb.OrderResponse
	for _, o := range orders {
		orderResponse := &pb.OrderResponse{
			Id:         o.ID,
			Price:      float32(o.Price),
			Tax:        float32(o.Tax),
			FinalPrice: float32(o.FinalPrice),
		}

		ordersResponse = append(ordersResponse, orderResponse)
	}

	return &pb.OrderListResponse{Itens: ordersResponse}, nil
}
