package services

import (
	"fmt"
	"github/rizkypedia/quiz-master/internal/database"
	"github/rizkypedia/quiz-master/internal/models"
	"strings"

	"github/rizkypedia/quiz-master/internal/quizrules"
	"gorm.io/gorm"
)

var ErrInvalidQuestion = quizrules.ErrInvalidQuestion

func GetQuestions() ([]models.Question, error) {
	var questions []models.Question

	err := database.DB.Preload("Answers", func(db *gorm.DB) *gorm.DB {
		return db.Order("answers.id ASC")
	}).
		Preload("Category").
		Preload("QuizType").
		Order("questions.id ASC").
		Find(&questions).Error
	if err != nil {
		return nil, err
	}

	return questions, nil
}

func CreateQuestion(question *models.Question) error {
	if question == nil || strings.TrimSpace(question.Question) == "" ||
		question.CategoryId == 0 || question.QuizTypeID == 0 {
		return fmt.Errorf("%w: question text, category, and quiz type are required", ErrInvalidQuestion)
	}
	if err := quizrules.ValidateAnswers(question.Answers); err != nil {
		return err
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		var category models.Category
		if err := tx.First(&category, question.CategoryId).Error; err != nil {
			return err
		}
		var quizType models.QuizType
		if err := tx.First(&quizType, question.QuizTypeID).Error; err != nil {
			return err
		}
		if err := quizrules.ValidateCorrectAnswerCount(quizType.Name, question.Answers); err != nil {
			return err
		}
		question.Question = strings.TrimSpace(question.Question)
		return tx.Create(question).Error
	})
}
