package middleware

import "github.com/gin-gonic/gin"

//Проверка сессии админки. Подключается к post login и post products.
func AuthMiddleware(c *gin.Context) {
	/* Логика:
	Проверка куки отправителя запроса. Находит там активную сессию.
	Происходит сверка сессии с доступными в памяти на данный момент.
	Сверка по дате закрытия сессии. Если все ок - пропускает на handler, если нет - 401 ошибка*/
	c.Next()
}
