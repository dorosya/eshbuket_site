package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func backendURL() string {
	url := os.Getenv("BACKEND_URL")
	if url == "" {
		return "http://localhost:8080"
	}
	return url
}

func TestGinMode(t *testing.T) {
	if gin.Mode() != gin.TestMode {
		t.Fatal("Gin is not in test mode")
	}
}

func TestCreateProductWithCategory(t *testing.T) {
	body := map[string]string{
		"name":     "Букет1",
		"price":    "100",
		"category": "flowers",
	}

	jsonBody, _ := json.Marshal(body)

	resp, err := http.Post(
		backendURL()+"/api/products",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
}

func TestCreateProductWithoutCategory(t *testing.T) {
	body := map[string]string{
		"name":  "Tulip",
		"price": "80",
	}

	jsonBody, _ := json.Marshal(body)

	resp, err := http.Post(
		backendURL()+"/api/products",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
}

func TestCreateProductInvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp, err := http.Post(
		backendURL()+"/api/products",
		"application/json",
		bytes.NewBufferString(`{"name":"Rose"}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}
