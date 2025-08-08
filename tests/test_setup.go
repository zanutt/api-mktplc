package tests

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/zanutt/api-mktplc/internal/database"
	"github.com/zanutt/api-mktplc/internal/product"
	"github.com/zanutt/api-mktplc/internal/router"
	"github.com/zanutt/api-mktplc/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var once sync.Once
var TestDB *gorm.DB
var TestRouter *gin.Engine

func SetupTestEnv() {
	once.Do(func() {
		_ = godotenv.Load("../.env.test")

		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
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

		database.RunMigrations(db)

		TestDB = db
		TestRouter = router.SetupRouter(db)
	})
}

// ResetDatabase remove todas as tabelas e recria do zero
func ResetDatabase() {
	err := TestDB.Migrator().DropTable(&user.User{}, &product.Product{})
	if err != nil {
		log.Fatalf("Erro ao dropar tabelas: %v", err)
	}

	database.RunMigrations(TestDB)
}
