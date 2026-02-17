package services

import (
	"errors"
	"github/rizkypedia/quiz-master/internal/database"
	"github/rizkypedia/quiz-master/internal/dto"
	"github/rizkypedia/quiz-master/internal/models"

	"gorm.io/gorm"
)

func GetCategories() ([]models.Category, error) {
	var categories []models.Category
	err := database.DB.Preload("QuizType").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func GetCategoryById(id uint) (*models.Category, error) {
	var cat models.Category
	err := database.DB.
		Preload("QuizType").
		First(&cat, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}
	return &cat, err
}

func CreateCategory(input dto.CreateCategoryDTO) error {
	category := models.Category{
		Name:       input.Name,
		QuizTypeID: input.QuizType.ID,
	}
	return database.DB.Create(&category).Error
}
