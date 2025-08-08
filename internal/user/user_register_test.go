package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing" // você esqueceu de importar "testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	_ = godotenv.Load("../../.env.test") // ajuste o caminho conforme necessário

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
		panic("failed to connect to test database: " + err.Error())
	}

	// Limpa e migra a tabela User
	_ = db.Migrator().DropTable(&User{})
	_ = db.AutoMigrate(&User{})

	return db
}

func TestRegisterHandler_Success(t *testing.T) {
	db := setupTestDB()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/register", RegisterHandler(db))

	payload := map[string]interface{}{
		"username": "joaozinho", // Adicione o campo username se for obrigatório
		"email":    "joao@example.com",
		"password": "senha123",
		"type":     "user",
	}

	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}
