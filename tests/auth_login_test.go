package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zanutt/api-mktplc/internal/user"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginHandler_Sucess(t *testing.T) {
	SetupTest(t)

	// Create a test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)

	u := user.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Type:     "user",
	}

	if err := TestDB.Create(&u).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Payload for login
	payload := map[string]string{
		"email":    "test@example.com",
		"password": "senha123",
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	// Perform the request
	TestRouter.ServeHTTP(res, req)

	// Validate Response
	if res.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", res.Code)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(res.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	if _, ok := result["token"]; !ok {
		t.Fatalf("Expected token in response")
	}
}

func TestLoginHandler_WrongPass(t *testing.T) {
	SetupTest(t)

	// Create a test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)

	u := user.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Type:     "user",
	}
	TestDB.Create(&u)

	// Payload for login with wrong password
	payload := map[string]string{
		"email":    "test@example.com",
		"password": "wrongpassword",
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	TestRouter.ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status 401, got %d", res.Code)
	}
}

func TestLoginHandler_NoExistingUser(t *testing.T) {
	SetupTest(t)

	payload := map[string]string{
		"email":    "nonecxiste@example.com",
		"password": "senha123",
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	TestRouter.ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status 401, got %d", res.Code)
	}
}
