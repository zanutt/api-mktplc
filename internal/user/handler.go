package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			Type     string `json:"type"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
			return
		}

		if input.Type != "admin" && input.Type != "user" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo de usuário inválido"})
			return
		}

		if input.Email == "" || input.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email e senha são obrigatórios"})
			return
		}

		hashed, err := HashPassword(input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar senha"})
			return
		}

		user := User{
			Email:    input.Email,
			Password: hashed,
			Type:     input.Type,
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar usuário"})
			return
		}

		token, err := GenerateJWT(user.ID, user.Type)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"token": token,
			"user": gin.H{
				"id":    user.ID,
				"email": user.Email,
				"type":  user.Type,
			},
		})
	}
}

func LoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
			return
		}

		if input.Email == "" || input.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email e senha são obrigatórios"})
			return
		}

		var user User
		if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar usuário"})
			}
			return
		}

		if err := CheckPasswordHash(input.Password, user.Password); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
			return
		}

		token, err := GenerateJWT(user.ID, user.Type)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"user": gin.H{
				"id":    user.ID,
				"email": user.Email,
				"type":  user.Type,
			},
		})

	}
}
