package services

import (
	"eshbuket/models"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Authenticate(login, password string) bool
	CreateSession(req models.LoginRequest) string
}

type authService struct{}

func NewAuthService() AuthService {
	return &authService{}
}

func (s *authService) Authenticate(login string, password string) bool {
	Storedhash := os.Getenv("ADMIN_PASSWORD_HASH")
	Adminlogin := os.Getenv("ADMIN_LOGIN")

	if login != Adminlogin {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(Storedhash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func (s *authService) CreateSession(req models.LoginRequest) string {
	var sessionID = uuid.NewString()
	models.Sessions[sessionID] = models.Session{
		Username: req.Login,
		Expires:  time.Now().Add(1 * time.Hour),
	}
	return sessionID
}
