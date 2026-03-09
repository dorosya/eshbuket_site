package auth

import (
	"eshbuket/internal/Domain/models"
	store "eshbuket/internal/repository/Store"
	"os"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Authenticate(t *testing.T) {
	t.Setenv("ADMIN_LOGIN", "boss")

	hash, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to generate bcrypt hash: %v", err)
	}
	t.Setenv("ADMIN_PASSWORD_HASH", string(hash))

	svc := NewAuthService(store.NewSessionStore())

	if !svc.Authenticate("boss", "secret123") {
		t.Fatal("expected successful authentication")
	}
	if svc.Authenticate("boss", "wrong") {
		t.Fatal("expected auth to fail with wrong password")
	}
	if svc.Authenticate("other", "secret123") {
		t.Fatal("expected auth to fail with wrong login")
	}
}

func TestAuthService_CreateAndValidateSession(t *testing.T) {
	st := store.NewSessionStore()
	svc := NewAuthService(st)

	sessionID, err := svc.CreateSession("boss")
	if err != nil {
		t.Fatalf("CreateSession returned error: %v", err)
	}
	if sessionID == "" {
		t.Fatal("expected non-empty session id")
	}
	if !svc.ValidateSession(sessionID) {
		t.Fatal("expected created session to be valid")
	}
}

func TestAuthService_ValidateSession_Expired(t *testing.T) {
	st := store.NewSessionStore()
	svc := NewAuthService(st)

	st.Set("old-session", models.Session{
		Username: "boss",
		Expires:  time.Now().Add(-1 * time.Second),
	})

	if svc.ValidateSession("old-session") {
		t.Fatal("expected expired session to be invalid")
	}
}

func TestAuthService_Authenticate_EmptyEnv(t *testing.T) {
	_ = os.Unsetenv("ADMIN_LOGIN")
	_ = os.Unsetenv("ADMIN_PASSWORD_HASH")
	svc := NewAuthService(store.NewSessionStore())

	if svc.Authenticate("any", "any") {
		t.Fatal("expected auth to fail when env vars are missing")
	}
}
