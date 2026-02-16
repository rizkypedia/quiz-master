package services

import (
	"errors"
	"github/rizkypedia/quiz-master/internal/database"
	"github/rizkypedia/quiz-master/internal/models"
	"strings"

	"gorm.io/gorm"
)

func GetQuestions() ([]models.Question, error) {
	var questions []models.Question

	// err := database.DB.Preload("Answers").Find(&questions).Error
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

func CreateQuestion(q string) error {
	qClean := strings.Trim(q, " ")
	length := len(qClean)
	if length == 0 {
		return errors.New("Empty question is not allowed")
	}
	question := models.Question{
		Question: qClean,
	}
	return database.DB.Create(&question).Error
}
