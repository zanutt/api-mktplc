package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zanutt/api-mktplc/internal/middleware"
	"github.com/zanutt/api-mktplc/internal/product"
	"github.com/zanutt/api-mktplc/internal/user"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Rotas públicas de autenticação
	r.POST("/register", user.RegisterHandler(db))
	r.POST("/login", user.LoginHandler(db))

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	auth.GET("/me", user.MeHandler(db))

	// Rotas protegidas
	productGroup := r.Group("/products")
	productGroup.GET("", product.NewProductHandler(db).List)

	// Admin-only
	adminGroup := productGroup.Group("")
	adminGroup.Use(middleware.AuthMiddleware("admin"))
	adminGroup.POST("", product.NewProductHandler(db).Create)
	adminGroup.PUT("/:id", product.NewProductHandler(db).Update)
	adminGroup.DELETE("/:id", product.NewProductHandler(db).Delete)

	return r
}
