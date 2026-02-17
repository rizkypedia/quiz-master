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
		dtos[i] = dto.CategoryDTO{ID: c.ID, Name: c.Name, QuizType: dto.QuizTypeDTO{ID: c.QuizType.ID, Name: c.QuizType.Name}}
	}

	c.JSON(http.StatusOK, gin.H{
		"items": dtos,
		"count": len(dtos),
	})
}

func GetCategoryById(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid id",
		})
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
	dtoObj := dto.CategoryDTO{
		ID:   cat.ID,
		Name: cat.Name,
		QuizType: dto.QuizTypeDTO{
			ID:   cat.QuizType.ID,
			Name: cat.QuizType.Name,
		},
	}
	c.JSON(http.StatusOK, dtoObj)
}
