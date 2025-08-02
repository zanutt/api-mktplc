package product

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	handler := NewProductHandler(db)

	r.GET("/products", handler.List)
	r.POST("/products", handler.Create)
	r.PUT("/products/:id", handler.Update)
	r.DELETE("/products/:id", handler.Delete)

	return r
}
