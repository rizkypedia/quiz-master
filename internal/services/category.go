package services

import (
	"errors"
	"strings"

	"github/rizkypedia/quiz-master/internal/database"
	"github/rizkypedia/quiz-master/internal/dto"
	"github/rizkypedia/quiz-master/internal/models"

	"gorm.io/gorm"
)

var ErrInvalidCategory = errors.New("invalid category")

func GetCategories() ([]models.Category, error) {
	var categories []models.Category
	err := database.DB.Order("categories.id ASC").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func GetCategoryById(id uint) (*models.Category, error) {
	var cat models.Category
	err := database.DB.First(&cat, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}
	return &cat, err
}

func CreateCategory(input dto.CreateCategoryDTO) error {
	if strings.TrimSpace(input.Name) == "" {
		return ErrInvalidCategory
	}
	category := models.Category{
		Name: strings.TrimSpace(input.Name),
	}
	return database.DB.Create(&category).Error
}
