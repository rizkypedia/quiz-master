package handlers

import (
	"net/http"

	"github/rizkypedia/quiz-master/internal/dto"
	"github/rizkypedia/quiz-master/internal/services"

	"github.com/gin-gonic/gin"
)

func GetQuizTypes(c *gin.Context) {
	quizTypes, err := services.GetQuizTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	items := make([]dto.QuizTypeDTO, len(quizTypes))
	for i, quizType := range quizTypes {
		items[i] = dto.QuizTypeDTO{ID: quizType.ID, Name: quizType.Name}
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"count": len(items),
	})
}
