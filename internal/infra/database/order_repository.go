package database

import (
	"database/sql"

	"github.com/alexsuriano/clean-architecture/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) FindAll() ([]entity.Order, error) {
	query := "SELECT id, price, tax, final_price FROM orders"
	rows, err := r.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order

	for rows.Next() {
		var id string
		var price, tax, finalPrice float64

		err := rows.Scan(&id, &price, &tax, &finalPrice)
		if err != nil {
			return nil, err
		}

		order := entity.Order{
			ID:         id,
			Price:      price,
			Tax:        tax,
			FinalPrice: finalPrice,
		}

		orders = append(orders, order)
	}

	return orders, nil
}
