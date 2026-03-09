package postgres

import (
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestBuildDSN(t *testing.T) {
	got := BuildDSN("localhost", "5432", "user", "pass", "shop")
	want := "host=localhost port=5432 user=user password=pass dbname=shop sslmode=disable"

	if got != want {
		t.Fatalf("unexpected dsn: got %q want %q", got, want)
	}
}

func TestRunMigrations_ErrorWhenDirDoesNotExist(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() { _ = db.Close() }()

	err = RunMigrations(db, "./migrations-do-not-exist")
	if err == nil {
		t.Fatal("expected error for missing migrations dir")
	}
	if !strings.Contains(err.Error(), "run migrations:") {
		t.Fatalf("unexpected error: %v", err)
	}
}

