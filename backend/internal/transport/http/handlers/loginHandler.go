package handlers

import (
	"eshbuket/internal/service/auth"
	"eshbuket/internal/transport/http/dto"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	service auth.AuthService
}

func NewLoginHandler(ls auth.AuthService) *LoginHandler {
	return &LoginHandler{ls}
}

// POST /api/login
func (h *LoginHandler) LoginHandler(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if !h.service.Authenticate(req.Login, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid login and/or password"})
		return
	}

	sessionID, err := h.service.CreateSession(req.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}

	c.SetSameSite(cookieSameSiteMode())
	c.SetCookie(
		"session_id",
		sessionID,
		3600,
		"/",
		"",
		shouldUseSecureCookie(),
		true,
	)
	c.JSON(http.StatusOK, gin.H{"message": "authorized"})
}

func shouldUseSecureCookie() bool {
	if raw := strings.TrimSpace(os.Getenv("COOKIE_SECURE")); raw != "" {
		if v, err := strconv.ParseBool(raw); err == nil {
			return v
		}
	}

	env := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	switch env {
	case "local", "dev", "development", "test":
		return false
	default:
		return true
	}
}

func cookieSameSiteMode() http.SameSite {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("COOKIE_SAMESITE"))) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}
