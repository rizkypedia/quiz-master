package services

import (
	"github/rizkypedia/quiz-master/internal/database"
	"github/rizkypedia/quiz-master/internal/models"

	"gorm.io/gorm"
)

func GetQuestions() ([]models.Question, error) {
	var questions []models.Question

	err := database.DB.Preload("Answers", func(db *gorm.DB) *gorm.DB {
		return db.Order("answers.id ASC")
	}).
		Preload("Category").
		Preload("Category.QuizType").
		Order("questions.id ASC").
		Find(&questions).Error
	if err != nil {
		return nil, err
	}

	return questions, err
}

func CreateQuestion(question models.Question) error {
	err := database.DB.Create(&question).Error
	if err != nil {
		return err
	}
	return nil
}
