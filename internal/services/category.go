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
	err := database.DB.Where("categories.active = ?", true).Order("categories.id ASC").Find(&categories).Error
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
		Name:   strings.TrimSpace(input.Name),
		Active: true,
	}
	return database.DB.Create(&category).Error
}

func GetActiveCategoryByID(id uint) (*models.Category, error) {
	var category models.Category
	err := database.DB.Where("active = ?", true).First(&category, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}
	return &category, err
}

func SetCategoryActive(id uint, active bool) error {
	result := database.DB.Model(&models.Category{}).Where("id = ?", id).Update("active", active)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
