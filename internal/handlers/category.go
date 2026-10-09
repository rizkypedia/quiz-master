package handlers

import (
	"errors"
	"github/rizkypedia/quiz-master/internal/dto"
	"github/rizkypedia/quiz-master/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetCategories(c *gin.Context) {
	cs, err := services.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	dtos := make([]dto.CategoryDTO, len(cs))
	for i, c := range cs {
		dtos[i] = dto.CategoryDTO{ID: c.ID, Name: c.Name, Active: c.Active}
	}

	c.JSON(http.StatusOK, gin.H{
		"items": dtos,
		"count": len(dtos),
	})
}

func GetCategoryById(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, strconv.IntSize)
	if err != nil || id64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid id",
		})
		return
	}
	id := uint(id64)

	cat, err := services.GetCategoryById(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	if !cat.Active {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	dtoObj := dto.CategoryDTO{
		ID:     cat.ID,
		Name:   cat.Name,
		Active: cat.Active,
	}
	c.JSON(http.StatusOK, dtoObj)
}

func SetCategoryActive(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, strconv.IntSize)
	if err != nil || id64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	var requestBody dto.SetCategoryActiveDTO
	if err := c.ShouldBindJSON(&requestBody); err != nil || requestBody.Active == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "active is required"})
		return
	}

	err = services.SetCategoryActive(uint(id64), *requestBody.Active)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "category not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "category updated"})
}

func CreateCategory(c *gin.Context) {
	var requestBody dto.CreateCategoryDTO
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := services.CreateCategory(requestBody); err != nil {
		if errors.Is(err, services.ErrInvalidCategory) {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "category created"})

}
