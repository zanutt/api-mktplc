package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterHandler(t *testing.T) {
	SetupTest(t)

	if TestDB == nil || TestRouter == nil {
		t.Fatal("TestDB ou TestRouter não foi inicializado. Verifique se SetupTestEnv está funcionando corretamente.")
	}

	gin.SetMode(gin.TestMode)

	payload := map[string]interface{}{
		"username": "adminuser",
		"email":    "test@example.com",
		"password": "123456",
		"type":     "admin",
	}

	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	TestRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}
