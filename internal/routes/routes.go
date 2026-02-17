package routes

import (
	"github/rizkypedia/quiz-master/internal/handlers"
	"github/rizkypedia/quiz-master/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.GET("/question", handlers.GetQuestions)
	r.POST("/question", handlers.AddQuestion)
	r.GET("/categories", handlers.GetCategories)
	r.GET("/category/:id", handlers.GetCategoryById)
	auth := r.Group("/")
	auth.Use(middleware.AuthRequired())
	auth.GET("/profile", handlers.Profile)
}
