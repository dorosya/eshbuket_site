package product

import (
	"context"
	"database/sql"
	"errors"
	"eshbuket/internal/Domain/models"
	"eshbuket/internal/transport/http/dto"
	"fmt"
	"log"
	"strconv"
)

type ProductService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{repo}
}

func (service *ProductService) ProductGetService(ctx context.Context, category string) ([]models.Product, error) {
	rows, err := service.repo.FindProducts(ctx, category)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("Failed to close rows:", err)
		}
	}()

	products := []models.Product{}
	for rows.Next() {
		var p models.Product
		var price, id sql.NullInt64
		var imagePath sql.NullString
		if err := rows.Scan(&id, &p.Name, &price, &p.Category, &imagePath); err != nil {
			return nil, err
		}
		p.ID = int(id.Int64)
		//лучше переделать под Decimal
		p.Price = int(price.Int64)
		if imagePath.Valid && imagePath.String != "" {
			p.ImageURL = "/api/products/" + fmt.Sprint(p.ID) + "/image"
		} else {
			p.ImageURL = ""
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (service *ProductService) ProductPostService(ctx context.Context, req dto.ProductRequest) error {
	price, err := strconv.Atoi(req.Price)
	if err != nil {
		return errors.New("invalid price format")
	}
	if price <= 0 {
		return errors.New("price must be positive")
	}
	err = service.repo.InsertProduct(ctx, req.Name, req.Price, req.Category)
	if err != nil {
		return err
	}
	return nil
}
