package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAndListProducts(t *testing.T) {
	SetupTest(t)

	// Registra usuário admin
	adminPayload := map[string]interface{}{
		"username": "adminuser",
		"email":    "admin@example.com",
		"password": "123456",
		"type":     "admin",
	}
	adminBody, _ := json.Marshal(adminPayload)
	req := httptest.NewRequest("POST", "/register", bytes.NewReader(adminBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	token := resp["token"]
	assert.NotEmpty(t, token)

	// Cria um produto
	productPayload := map[string]interface{}{
		"name":        "Smartphone",
		"description": "Android 12, 128GB",
		"price":       2500.00,
		"category":    "Eletrônicos",
	}
	productBody, _ := json.Marshal(productPayload)
	req = httptest.NewRequest("POST", "/products", bytes.NewReader(productBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Lista produtos
	req = httptest.NewRequest("GET", "/products", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var respList map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &respList)

	data, ok := respList["data"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, data, 1)
}
