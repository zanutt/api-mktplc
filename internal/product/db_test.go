package product

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	// Carrega variáveis do .env.test
	cwd, _ := os.Getwd()
	envPath := cwd + "/../../.env.test"
	err := godotenv.Load(envPath)
	if err != nil {
		panic("could not load .env.test")
	}

	// Exibe os valores carregados (para debug)
	fmt.Println("host:", os.Getenv("DB_HOST"))
	fmt.Println("port:", os.Getenv("DB_PORT"))
	fmt.Println("user:", os.Getenv("DB_USER"))
	fmt.Println("password:", os.Getenv("DB_PASSWORD"))
	fmt.Println("dbname:", os.Getenv("DB_NAME"))

	// Monta a DSN
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	// Conecta ao banco
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
	}

	// Limpa e recria a tabela de produto
	_ = db.Migrator().DropTable(&Product{})
	_ = db.AutoMigrate(&Product{})

	// Atribui o banco para uso nos testes
	testDB = db

	// Executa os testes
	code := m.Run()
	os.Exit(code)
}
