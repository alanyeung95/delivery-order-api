package mysql

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/alanyeung95/delivery-order-api/pkg/orders"
)

// NewOrderRepository is the repo to store account model
func NewOrderRepository(db *sql.DB) (*OrderRepository, error) {
	return &OrderRepository{db}, nil
}

type OrderRepository struct {
	repository *sql.DB
}

// interface check
var _ orders.Repository = (*OrderRepository)(nil)

func (r *OrderRepository) FindByID(ctx context.Context, id string) (*orders.Order, error) {
	var order orders.Order
	// TODO: db naming
	err := r.repository.QueryRow("SELECT id, distance, status FROM delivery_order where id = ?", id).Scan(&order.ID, &order.Distance, &order.Status)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) Find(ctx context.Context, limit int, offset int) ([]orders.Order, error) {
	// Execute the query
	results, err := r.repository.Query("SELECT * FROM delivery_order LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}
	var orderList []orders.Order

	for results.Next() {
		var order orders.Order
		// for each row, scan the result into our tag composite object
		err = results.Scan(&order.ID, &order.Distance, &order.Status)
		if err != nil {
			return nil, err
		}
		orderList = append(orderList, order)
	}

	if orderList == nil {
		return []orders.Order{}, nil
	}

	return orderList, nil
}

func (r *OrderRepository) Create(ctx context.Context, id string, distance int, status string) error {
	_, err := r.repository.Exec("INSERT INTO delivery_order(id,distance,status) values(?,?,?)", id, distance, status)
	return err
}

func (r *OrderRepository) Update(ctx context.Context, id string) error {
	result, err := r.repository.Exec("Update delivery_order set status=? where id=? and status=?", "TAKEN", id, "UNASSIGNED")
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if count != 1 || err != nil {
		return err
	}

	return nil
}
