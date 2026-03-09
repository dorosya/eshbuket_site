package order

import (
	"context"
	"database/sql"
)

type OrderRepository interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)

	GetProductPrice(
		ctx context.Context,
		tx *sql.Tx,
		productID int,
	) (int64, error)

	CreateOrder(
		ctx context.Context,
		tx *sql.Tx,
		contactData string,
		comment string,
		totalPriceRub int64,
	) (int, error)

	AddOrderProduct(
		ctx context.Context,
		tx *sql.Tx,
		orderID int,
		productID int,
		quantity int,
	) error
}
