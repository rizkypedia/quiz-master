package handlers

import (
	"errors"
	"github/rizkypedia/quiz-master/internal/dto"
	"github/rizkypedia/quiz-master/internal/models"
	"github/rizkypedia/quiz-master/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
			QuizTypeID: q.QuizTypeID,
			QuizType:   dto.QuizTypeDTO{ID: q.QuizType.ID, Name: q.QuizType.Name},
			Category: dto.CategoryDTO{
				ID: q.Category.ID, Name: q.Category.Name,
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

	q, err := prepareQuestion(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if err := services.CreateQuestion(q); err != nil {
		if errors.Is(err, services.ErrInvalidQuestion) {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "category or quiz type not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": "created",
	})
}

func prepareQuestion(input dto.CreateQuestionDTO) (*models.Question, error) {
	if input.Question == nil || input.CategoryId == nil || input.QuizTypeID == nil || input.Answers == nil {
		return nil, errors.New("question, categoryId, quizTypeId, and answers are required")
	}

	q := models.Question{
		Question:   *input.Question,
		CategoryId: *input.CategoryId,
		QuizTypeID: *input.QuizTypeID,
	}

	for _, a := range *input.Answers {
		if a.Answer == nil || a.IsCorrect == nil {
			return nil, errors.New("each answer requires answer text and isCorrect")
		}
		q.Answers = append(q.Answers, models.Answers{
			Answer:    *a.Answer,
			IsCorrect: *a.IsCorrect,
		})
	}

	return &q, nil
}
