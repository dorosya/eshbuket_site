package handlers

import (
	"bytes"
	store "eshbuket/internal/repository/Store"
	"eshbuket/internal/service/auth"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginHandler_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewLoginHandler(*auth.NewAuthService(store.NewSessionStore()))

	router := gin.New()
	router.POST("/api/login", handler.LoginHandler)

	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString(`{"login":"admin"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestLoginHandler_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("ADMIN_LOGIN", "admin")
	hash, err := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	t.Setenv("ADMIN_PASSWORD_HASH", string(hash))

	handler := NewLoginHandler(*auth.NewAuthService(store.NewSessionStore()))
	router := gin.New()
	router.POST("/api/login", handler.LoginHandler)

	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString(`{"login":"admin","password":"wrong"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
	if strings.Contains(strings.Join(rec.Header().Values("Set-Cookie"), ";"), "session_id=") {
		t.Fatal("did not expect session cookie for unauthorized login")
	}
}

func TestLoginHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("ADMIN_LOGIN", "admin")
	hash, err := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	t.Setenv("ADMIN_PASSWORD_HASH", string(hash))

	handler := NewLoginHandler(*auth.NewAuthService(store.NewSessionStore()))
	router := gin.New()
	router.POST("/api/login", handler.LoginHandler)

	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString(`{"login":"admin","password":"correct"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(strings.Join(rec.Header().Values("Set-Cookie"), ";"), "session_id=") {
		t.Fatal("expected session cookie in response")
	}
}

func TestShouldUseSecureCookie(t *testing.T) {
	t.Run("env local uses insecure", func(t *testing.T) {
		t.Setenv("COOKIE_SECURE", "")
		t.Setenv("APP_ENV", "local")
		if shouldUseSecureCookie() {
			t.Fatal("expected insecure cookie for local env")
		}
	})

	t.Run("explicit override wins", func(t *testing.T) {
		t.Setenv("APP_ENV", "local")
		t.Setenv("COOKIE_SECURE", "true")
		if !shouldUseSecureCookie() {
			t.Fatal("expected secure cookie when COOKIE_SECURE=true")
		}
	})
}
