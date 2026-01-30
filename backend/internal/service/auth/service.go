package auth

import (
	"eshbuket/internal/Domain/models"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Authenticate(login string, password string) bool {
	Storedhash := os.Getenv("ADMIN_PASSWORD_HASH")
	Adminlogin := os.Getenv("ADMIN_LOGIN")

	if login != Adminlogin {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(Storedhash), []byte(password))
	return err == nil
}

// TODO: JWT AUTH
func CreateSession(Login string) string {
	var sessionID = uuid.NewString()
	models.Sessions[sessionID] = models.Session{
		Username: Login,
		Expires:  time.Now().Add(1 * time.Hour),
	}
	return sessionID
}
