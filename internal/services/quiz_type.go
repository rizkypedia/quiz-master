package services

import (
	"github/rizkypedia/quiz-master/internal/database"
	"github/rizkypedia/quiz-master/internal/models"
)

func GetQuizTypes() ([]models.QuizType, error) {
	var quizTypes []models.QuizType
	err := database.DB.Order("quiz_types.id ASC").Find(&quizTypes).Error
	if err != nil {
		return nil, err
	}
	return quizTypes, nil
}
