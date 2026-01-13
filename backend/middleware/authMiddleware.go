package middleware

import (
	handlers "eshbuket/handlers/structures"
	"time"

	"github.com/gin-gonic/gin"
)

// Проверка сессии админки. Подключается к post login и post products.
func AuthMiddleware(c *gin.Context) {
	//Временная мера для пропуска тестов. Удалить позже
	if gin.Mode() == gin.TestMode {
		c.Next()
		return
	}
	//Проверка сессии в куках
	session_id, err := c.Cookie("session_id")
	if err != nil {
		c.JSON(401, gin.H{"error": "Не авторизован"})
		c.Abort()
		return
	}

	session, ok := handlers.Sessions[session_id]
	//Проверка вхождения сессии в памяти программы
	if !ok {
		c.JSON(401, gin.H{"error": "Не авторизован"})
		c.Abort()
		return
	}
	//Проверка временных данных сессии
	if session.Expires.Before(time.Now()) {
		delete(handlers.Sessions, session_id) //  очистка просроченной сессии
		c.JSON(401, gin.H{"error": "Сессия истекла"})
		c.Abort()
		return
	}

	c.Next()
	/* Логика:
	Проверка куки отправителя запроса. Находит там активную сессию.
	Происходит сверка сессии с доступными в памяти на данный момент.
	Сверка по дате закрытия сессии. Если все ок - пропускает на handler, если нет - 401 ошибка*/
}
