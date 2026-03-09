package order

import (
	"context"
	"eshbuket/internal/transport/http/dto"
)

type OrderService struct {
	repo OrderRepository
}

func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{repo}
}

func (service *OrderService) OrderServiceFunc(ctx context.Context, req dto.OrderRequest) (orderID int, totalPriceRub int64, err error) {
	tx, err := service.repo.BeginTx(ctx)
	if err != nil {
		return 0, 0, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	for _, p := range req.Products {
		var priceRub int64
		priceRub, err = service.repo.GetProductPrice(ctx, tx, p.ProductID)
		if err != nil {
			return
		}
		totalPriceRub += priceRub * int64(p.Quantity)
	}

	orderID, err = service.repo.CreateOrder(ctx, tx, req.ContactData, req.Comment, totalPriceRub)
	if err != nil {
		return
	}

	for _, product := range req.Products {
		err = service.repo.AddOrderProduct(ctx, tx, orderID, product.ProductID, product.Quantity)
		if err != nil {
			return
		}
	}

	if err = tx.Commit(); err != nil {
		return
	}

	return orderID, totalPriceRub, nil
}

