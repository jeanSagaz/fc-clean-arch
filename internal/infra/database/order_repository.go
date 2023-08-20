package database

import (
	"database/sql"

	"github.com/jeanSagaz/fc-clean-arch/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) GetAll() ([]entity.Order, error) {
	rows, err := r.Db.Query("Select id, price, tax, final_price from orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var order entity.Order
		err = rows.Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepository) GetById(id string) (*entity.Order, error) {
	stmt, err := r.Db.Prepare("Select id, price, tax, final_price FROM orders WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var order entity.Order
	err = stmt.QueryRow(id).Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) GetTotal() (int, error) {
	var total int
	err := r.Db.QueryRow("Select count(*) from orders").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) Update(order *entity.Order) error {
	stmt, err := r.Db.Prepare("UPDATE orders SET price = ?, tax = ?, final_price = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(order.Price, order.Tax, order.FinalPrice, order.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) Delete(id string) error {
	stmt, err := r.Db.Prepare("DELETE FROM orders WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
