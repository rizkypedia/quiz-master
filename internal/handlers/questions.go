package handlers

import (
	"github/rizkypedia/quiz-master/internal/dto"
	"github/rizkypedia/quiz-master/internal/models"
	"github/rizkypedia/quiz-master/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetQuestions(c *gin.Context) {
	qs, err := services.GetQuestions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	dtos := make([]dto.QuestionDTO, len(qs))
	for i, q := range qs {
		answers := make([]dto.AnswerDTO, len(q.Answers))
		for j, a := range q.Answers {
			answers[j] = dto.AnswerDTO{
				ID: a.ID, Answer: a.Answer, QuestionId: a.QuestionId, IsCorrect: a.IsCorrect,
			}
		}
		dtos[i] = dto.QuestionDTO{
			ID: q.ID, Question: q.Question, CategoryId: q.CategoryId,
			Category: dto.CategoryDTO{
				ID: q.Category.ID, Name: q.Category.Name,
				QuizType: dto.QuizTypeDTO{
					ID: q.Category.QuizType.ID, Name: q.Category.QuizType.Name,
				},
			},
			Answers: answers,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"items": dtos,
		"count": len(dtos),
	})
}

func CreateQuestion(c *gin.Context) {
	var input dto.CreateQuestionDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if input.Question == nil || input.CategoryId == nil || input.Answers == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing require fields",
		})
	}
	q := models.Question{
		Question:   *input.Question,
		CategoryId: *input.CategoryId,
	}

	if input.Answers != nil {
		for _, a := range *input.Answers {
			if a.Answer == nil || a.IsCorrect == nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "Answers field is required",
				})
				return
			}
			q.Answers = append(q.Answers, models.Answers{
				Answer:    *a.Answer,
				IsCorrect: *a.IsCorrect,
			})
		}
	}

	if err := services.CreateQuestion(q); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": "created",
	})
}
