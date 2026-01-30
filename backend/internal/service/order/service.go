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

func (service *OrderService) OrderServiceFunc(ctx context.Context, req dto.OrderRequest) (orderid int, totalprice int, err error) {
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
		var price int
		price, err = service.repo.GetProductPrice(ctx, tx, p.ProductID)
		if err != nil {
			return
		}

		totalprice += price * p.Quantity
	}

	// Вставка заказа
	orderid, err = service.repo.CreateOrder(ctx, tx, req.ContactData, req.Comment, totalprice)
	if err != nil {
		return
	}
	// Вставляем записи в order_products. Подтверждает транзакцию
	for _, product := range req.Products {
		err = service.repo.AddOrderProduct(ctx, tx, orderid, product.ProductID, product.Quantity)
		if err != nil {
			return
		}
	}
	if err = tx.Commit(); err != nil {
		return
	}
	return orderid, totalprice, nil
}
