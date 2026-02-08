package auth

import (
	"crypto/rand"
	"encoding/hex"
	"eshbuket/internal/Domain/models"
	store "eshbuket/internal/repository/Store"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	store *store.SessionStore
}

func NewAuthService(store *store.SessionStore) *AuthService {
	return &AuthService{store}
}
func (s *AuthService) CreateSession(username string) string {
	id := generateSessionID()
	s.store.Set(id, models.Session{
		Username: username,
		Expires:  time.Now().Add(1 * time.Hour),
	})
	return id
}

func (service *AuthService) Authenticate(login string, password string) bool {
	Storedhash := os.Getenv("ADMIN_PASSWORD_HASH")
	Adminlogin := os.Getenv("ADMIN_LOGIN")
	if login != Adminlogin {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(Storedhash), []byte(password))
	return err == nil
}

func (s *AuthService) ValidateSession(id string) bool {
	_, ok := s.store.Get(id)
	return ok
}

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
