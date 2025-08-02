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

	assert.Equal(t, http.StatusCreated, w.Code)

	// envia um GET
	req = httptest.NewRequest("GET", "/products", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var products []product.Product
	_ = json.Unmarshal(w.Body.Bytes(), &products)
	assert.Len(t, products, 1)
	assert.Equal(t, "Smartwatch", products[0].Name)
}
