package postgres

import (
	"context"
	"database/sql"
	"errors"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (repo *ProductRepository) FindProducts(ctx context.Context, category string) (*sql.Rows, error) {
	rows, err := repo.db.QueryContext(ctx, `
			SELECT id, name, price, category, image_path
			FROM products
			WHERE ($1 = '' OR category = $1)
		`, category)

	if err != nil {
		return rows, errors.New("repository error: cannot find an item during scanning products rows")
	}

	return rows, err
}

func (repo *ProductRepository) InsertProduct(ctx context.Context, name string, priceCents int64, category string) error {
	_, err := repo.db.ExecContext(ctx,
		"INSERT INTO products (name, price, category) VALUES ($1, $2, $3)",
		name, priceCents, category,
	)
	if err != nil {
		return errors.New("repository error: " + err.Error())
	}
	return nil
}
