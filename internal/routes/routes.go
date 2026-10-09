package routes

import (
	"github/rizkypedia/quiz-master/internal/handlers"
	"github/rizkypedia/quiz-master/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.GET("/question", handlers.GetQuizByToken)
	r.POST("/question", handlers.CreateQuestion)
	r.GET("/quiz", handlers.GetQuizByToken)
	r.POST("/quiz/result", handlers.SubmitQuizResult)
	r.GET("/quiz-link/validate", handlers.ValidateQuizLinkHandler)
	r.POST("/quiz-link", handlers.CreateQuizLink)
	r.GET("/categories", handlers.GetCategories)
	r.GET("/category/:id", handlers.GetCategoryById)
	r.POST("/category", handlers.CreateCategory)
	r.GET("/quiz-types", handlers.GetQuizTypes)
	auth := r.Group("/")
	auth.Use(middleware.AuthRequired())
	auth.GET("/profile", handlers.Profile)
	auth.PUT("/category/:id/active", handlers.SetCategoryActive)
}
