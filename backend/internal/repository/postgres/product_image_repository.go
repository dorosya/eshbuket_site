package postgres

import (
	"context"
	"database/sql"
)

type ProductImageRepository struct {
	db *sql.DB
}

func NewProductImageRepository(db *sql.DB) *ProductImageRepository {
	return &ProductImageRepository{db: db}
}

func (r *ProductImageRepository) ProductExists(ctx context.Context, productID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM products WHERE id=$1)", productID).Scan(&exists)
	return exists, err
}

func (r *ProductImageRepository) GetImagePath(ctx context.Context, productID string) (string, bool, error) {
	var p sql.NullString
	err := r.db.QueryRowContext(ctx, "SELECT image_path FROM products WHERE id=$1", productID).Scan(&p)
	if err != nil {
		return "", false, err
	}
	if !p.Valid || p.String == "" {
		return "", false, nil
	}
	return p.String, true, nil
}

func (r *ProductImageRepository) SetImagePath(ctx context.Context, productID string, relPath string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE products SET image_path=$1 WHERE id=$2", relPath, productID)
	return err
}
