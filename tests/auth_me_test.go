package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zanutt/api-mktplc/internal/user"
	"golang.org/x/crypto/bcrypt"
)

func TestMeHandler_Sucess(t *testing.T) {
	SetupTest(t)

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

	token, err := user.GenerateJWT(u.ID, u.Type)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	req, _ := http.NewRequest("GET", "/me", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	res := httptest.NewRecorder()

	TestRouter.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.Code)
	}
}

func TestMeHandler_Failure(t *testing.T) {
	SetupTest(t)

	req, _ := http.NewRequest("GET", "/me", nil)
	res := httptest.NewRecorder()

	TestRouter.ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, res.Code)
	}
}
