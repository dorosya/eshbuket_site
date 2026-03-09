package store

import (
	"eshbuket/internal/Domain/models"
	"testing"
	"time"
)

func TestSessionStore_SetAndGet(t *testing.T) {
	s := NewSessionStore()
	expected := models.Session{
		Username: "admin",
		Expires:  time.Now().Add(30 * time.Minute),
	}

	s.Set("sid-1", expected)
	got, ok := s.Get("sid-1")

	if !ok {
		t.Fatal("expected session to exist")
	}
	if got.Username != expected.Username {
		t.Fatalf("unexpected username: got %q want %q", got.Username, expected.Username)
	}
}

func TestSessionStore_GetExpiredSessionDeletesIt(t *testing.T) {
	s := NewSessionStore()
	s.Set("expired", models.Session{
		Username: "user",
		Expires:  time.Now().Add(-1 * time.Minute),
	})

	_, ok := s.Get("expired")
	if ok {
		t.Fatal("expected expired session to be invalid")
	}

	_, existsAfterDelete := s.sessions["expired"]
	if existsAfterDelete {
		t.Fatal("expected expired session to be deleted from store")
	}
}

func TestSessionStore_Delete(t *testing.T) {
	s := NewSessionStore()
	s.Set("sid-2", models.Session{
		Username: "user2",
		Expires:  time.Now().Add(10 * time.Minute),
	})

	s.Delete("sid-2")
	_, ok := s.Get("sid-2")
	if ok {
		t.Fatal("expected deleted session to be absent")
	}
}

