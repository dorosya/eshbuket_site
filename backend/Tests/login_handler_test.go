package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"eshbuket/handlers"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginHandler(t *testing.T) {
	// Настраиваем Gin в тестовом режиме
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Подготовка environment variables для теста
	adminPassword := "123"
	hash, _ := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	os.Setenv("ADMIN_PASSWORD_HASH", string(hash))
	os.Setenv("ADMIN_LOGIN", "admin")

	router.POST("/api/login", handlers.LoginHandler)

	// --- Тест: успешный логин ---
	reqBody := []byte(`{"login":"admin","password":"123"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// --- Тест: неверный пароль ---
	reqBody = []byte(`{"login":"admin","password":"wrongpass"}`)
	req = httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for wrong password, got %d", w.Code)
	}

	// --- Тест: неверный логин ---
	reqBody = []byte(`{"login":"wronglogin","password":"123"}`)
	req = httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for wrong login, got %d", w.Code)
	}
}
