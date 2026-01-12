package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	handlers "eshbuket/handlers/structures"
	"eshbuket/middleware"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware(t *testing.T) {
	// Включаем test mode
	gin.SetMode(gin.TestMode)

	// Создаём временную сессию
	sessionID := "test-session-123"
	handlers.Sessions = map[string]handlers.Session{
		sessionID: {
			Username: "admin",
			Expires:  time.Now().Add(1 * time.Hour),
		},
	}

	// Создаём роутер с middleware и пустым handler'ом
	router := gin.New()
	router.POST("/products", middleware.AuthMiddleware, func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 1️⃣ Тест без cookie → должен вернуть 401
	req1 := httptest.NewRequest("POST", "/products", nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w1.Code)
	}

	// 2️⃣ Тест с неправильным session_id → 401
	req2 := httptest.NewRequest("POST", "/products", nil)
	req2.AddCookie(&http.Cookie{Name: "session_id", Value: "wrong-session"})
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w2.Code)
	}

	// 3️⃣ Тест с правильным session_id → 200
	req3 := httptest.NewRequest("POST", "/products", nil)
	req3.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID})
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w3.Code)
	}
}
