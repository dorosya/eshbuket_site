package handlers

import (
	"eshbuket/models"
	"eshbuket/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type handler 

// POST /api/login - логин для админки
func LoginHandler(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	/* Логика реализации следующая:
	Окно ввода логина пароля. Отправляется пост запрос с содержанием header {
	 login: "admin";
	password: "123";
	}
	Оно сверяется с реальными данными, далее, если они корректны, создаёт новую сессию, которая сохраняется в памяти
	(в идеале redis, но думаю и в самой памяти тоже пойдет на время пока я не прикручу JWT), а так же в куки.
	Далее при каждом запросе на данный эндпоинт идет проверка через middleware, и, в случае если сессия активна, то пользователя пропускает в
	админ панель.*/

	auth := services.NewAuthService()
	if !auth.Authenticate(req.Login, req.Password) {
		c.JSON(401, gin.H{"error": "Неверно введен логин и/или пароль"})
		return
	}

	sessionID := auth.CreateSession(req)

	c.SetCookie(
		"session_id",
		sessionID,
		3600, // время жизни в секундах
		"/",
		"",   // домен
		true, // Secure (только HTTPS)
		true, // HttpOnly
	)
	c.JSON(http.StatusOK, gin.H{"message": "Авторизирован успешно"})
}
