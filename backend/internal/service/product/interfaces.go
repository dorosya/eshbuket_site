package product

import (
	"context"
	"database/sql"
)

type ProductRepository interface {
	FindProducts(ctx context.Context, category string) (*sql.Rows, error)
	InsertProduct(ctx context.Context, name string, priceRub int64, category string) error
}
