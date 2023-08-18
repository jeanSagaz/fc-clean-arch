package entity

type OrderRepositoryInterface interface {
	GetAll() ([]Order, error)
	GetById(id string) (*Order, error)
	GetTotal() (int, error)
	Save(order *Order) error
	Update(order *Order) error
	Delete(id string) error
}
