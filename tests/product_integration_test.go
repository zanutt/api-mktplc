package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/zanutt/api-mktplc/internal/product"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB
var router *gin.Engine

func TestMain(m *testing.M) {
	_ = godotenv.Load("../.env.test")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("cannot connect to test DB: " + err.Error())
	}

	_ = db.Migrator().DropTable(&product.Product{})
	_ = db.AutoMigrate(&product.Product{})

	testDB = db
	router = product.SetupRouter(db)

	os.Exit(m.Run())
}

func TestCreateAndListProducts(t *testing.T) {
	// limpa antes
	testDB.Exec("DELETE FROM products")

	// envia um POST
	payload := map[string]interface{}{
		"name":     "Smartwatch",
		"desc":     "GPS, AMOLED",
		"price":    499.90,
		"category": "Eletrônicos",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Println("POST /products response:", w.Body.String())

	assert.Equal(t, http.StatusCreated, w.Code)

	// envia um GET
	req = httptest.NewRequest("GET", "/products", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Println("GET /products response:", w.Body.String())

	assert.Equal(t, http.StatusOK, w.Code)

	type ProductsResponse struct {
		Data []product.Product `json:"data"`
		Meta struct {
			Limit         int `json:"limit"`
			Page          int `json:"page"`
			TotalPages    int `json:"totalPages"`
			TotalProducts int `json:"totalProducts"`
		} `json:"meta"`
	}

	var resp ProductsResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err, "Erro ao deserializar resposta do GET /products")
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, "Smartwatch", resp.Data[0].Name)
}
