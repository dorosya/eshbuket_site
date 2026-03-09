package middleware

import (
	store "eshbuket/internal/repository/Store"
	"eshbuket/internal/service/auth"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware_UnauthorizedWithoutCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := auth.NewAuthService(store.NewSessionStore())
	router := gin.New()
	router.Use(AuthMiddleware(svc))
	router.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestAuthMiddleware_AllowsValidSession(t *testing.T) {
	gin.SetMode(gin.TestMode)

	st := store.NewSessionStore()
	svc := auth.NewAuthService(st)
	sessionID, err := svc.CreateSession("admin")
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	router := gin.New()
	router.Use(AuthMiddleware(svc))
	router.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID})
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

