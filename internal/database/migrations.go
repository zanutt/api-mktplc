package database

import (
	"log"

	"github.com/zanutt/api-mktplc/internal/product"
	"github.com/zanutt/api-mktplc/internal/user"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	err := db.AutoMigrate(
		&user.User{},
		&product.Product{},
	)
	if err != nil {
		log.Fatalf("Erro ao executar migrações: %v", err)
	}
}
