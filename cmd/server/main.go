package main

import (
	"github/rizkypedia/quiz-master/internal/database"
	"github/rizkypedia/quiz-master/internal/routes"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	r := gin.Default()

	allowedOrigins := []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:3000"}
	if frontendURL := os.Getenv("FRONTEND_URL"); frontendURL != "" {
		allowedOrigins = append([]string{frontendURL}, allowedOrigins...)
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	routes.RegisterRoutes(r)

	r.Run(":8080")
}
