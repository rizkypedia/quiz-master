package handlers

import (
	"net/http"
	"os"
	"time"

	"github/rizkypedia/quiz-master/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Register(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	c.BindJSON(&body)

	if err := services.CreateUser(body.Email, body.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "created"})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	c.BindJSON(&body)

	user, ok := services.Authenticate(body.Email, body.Password)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Profile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "secure profile"})
}
