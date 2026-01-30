package postgres

import (
	"context"
	"database/sql"
	"errors"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (repo *OrderRepository) BeginTx() (*sql.Tx, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, errors.New("Repository error: failed to create transaction")
	}

	return tx, nil
}

func (repo *OrderRepository) GetProductPrice(ctx context.Context, tx *sql.Tx, productID int) (int, error) {
	var price int
	err := tx.QueryRowContext(ctx, "SELECT price FROM products WHERE id=$1", productID).Scan(&price)
	if err != nil {
		return 0, errors.New("Repository error: failed to find a product in order")
	}
	return price, nil
}

func (repo *OrderRepository) CreateOrder(ctx context.Context, tx *sql.Tx, contactData string, comment string, totalPrice int) (int, error) {
	var orderID int
	err := tx.QueryRowContext(ctx,
		"INSERT INTO orders (contact_data, comment, total_price) VALUES ($1, $2, $3) RETURNING id",
		contactData, comment, totalPrice,
	).Scan(&orderID)
	if err != nil {
		return 0, errors.New("Repository error: failed to insert Order")
	}
	return orderID, nil
}

func (repo *OrderRepository) AddOrderProduct(ctx context.Context, tx *sql.Tx, orderID int, productID int, quantity int) error {
	_, err := tx.ExecContext(ctx,
		"INSERT INTO order_products (order_id, product_id, quantity) VALUES ($1, $2, $3)",
		orderID, productID, quantity,
	)
	if err != nil {
		return errors.New("failed to insert order products")
	}
	return nil
}
