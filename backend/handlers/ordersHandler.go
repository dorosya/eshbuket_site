package handlers

import "github.com/gin-gonic/gin"

// POST /api/orders:
func OrdersHandler(c *gin.Context) {
	/* Отправляет на сервер следующие данные;
	Товар
	Контактные данные
	Комментарий к заказу(предпочтительный номер связи, дату и тд)
	Далее сервер в запросе обращается к БД, создает новый order с полученными данными. */
}
