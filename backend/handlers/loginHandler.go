package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// POST /api/login - логин для админки
func LoginHandler(c *gin.Context) {
	var req LoginRequest

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
	Storedhash := os.Getenv("ADMIN_PASSWORD_HASH")
	Adminlogin := os.Getenv("ADMIN_LOGIN")

	if req.Login != Adminlogin {
		c.JSON(401, "Неверно введен логин и/или пароль")
	}

	err := bcrypt.CompareHashAndPassword([]byte(Storedhash), []byte(req.Password))
	if err != nil {
		c.JSON(401, "Неверно введен пароль")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Авторизирован успешно"})
}
