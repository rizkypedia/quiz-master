package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github/rizkypedia/quiz-master/internal/dto"
	"github/rizkypedia/quiz-master/internal/models"
	"github/rizkypedia/quiz-master/internal/security"
	"github/rizkypedia/quiz-master/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateQuizLink(c *gin.Context) {
	var payload struct {
		CategoryID uint `json:"categoryId"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil || payload.CategoryID == 0 {
		categoryIDRaw := c.Query("categoryId")
		if categoryIDRaw == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "categoryId is required"})
			return
		}

		categoryID, err := strconv.ParseUint(categoryIDRaw, 10, 64)
		if err != nil || categoryID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "categoryId must be a valid positive integer"})
			return
		}
		payload.CategoryID = uint(categoryID)
	}

	if _, err := services.GetActiveCategoryByID(payload.CategoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "active category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tokenString, expiresAt, err := security.IssueCategoryQuizLink(payload.CategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	baseURL := strings.TrimRight(os.Getenv("FRONTEND_URL"), "/")
	if baseURL == "" {
		baseURL = "http://localhost:5173"
	}

	c.JSON(http.StatusOK, gin.H{
		"categoryId": payload.CategoryID,
		"token":      tokenString,
		"expiresAt":  expiresAt.Format(time.RFC3339),
		"link":       fmt.Sprintf("%s/?token=%s", baseURL, tokenString),
	})
}

func ValidateQuizLinkHandler(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "a secure quiz link is required"})
		return
	}

	categoryID, expiresAt, err := security.ValidateCategoryQuizLink(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "expiresAt": expiresAt.Format(time.RFC3339)})
		return
	}
	if _, err := services.GetActiveCategoryByID(categoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "this category is inactive"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categoryId": categoryID,
		"expiresAt":  expiresAt.Format(time.RFC3339),
		"valid":      true,
	})
}

func GetQuizByToken(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "a secure quiz link is required"})
		return
	}

	categoryID, expiresAt, err := security.ValidateCategoryQuizLink(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "expiresAt": expiresAt.Format(time.RFC3339)})
		return
	}
	if _, err := services.GetActiveCategoryByID(categoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "this category is inactive"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	questions, err := services.GetQuestions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filtered := make([]dto.QuestionDTO, 0)
	for _, q := range questions {
		if q.CategoryId != categoryID {
			continue
		}

		answers := make([]dto.AnswerDTO, 0, len(q.Answers))
		for _, a := range q.Answers {
			answers = append(answers, dto.AnswerDTO{ID: a.ID, Answer: a.Answer, QuestionId: a.QuestionId})
		}

		filtered = append(filtered, dto.QuestionDTO{
			ID:         q.ID,
			Question:   q.Question,
			CategoryId: q.CategoryId,
			QuizTypeID: q.QuizTypeID,
			QuizType:   dto.QuizTypeDTO{ID: q.QuizType.ID, Name: q.QuizType.Name},
			Category:   dto.CategoryDTO{ID: q.Category.ID, Name: q.Category.Name, Active: q.Category.Active},
			Answers:    answers,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"categoryId": categoryID,
		"expiresAt":  expiresAt.Format(time.RFC3339),
		"items":      filtered,
		"count":      len(filtered),
	})
}

func SubmitQuizResult(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "a secure quiz link is required"})
		return
	}

	categoryID, _, err := security.ValidateCategoryQuizLink(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if _, err := services.GetActiveCategoryByID(categoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "this category is inactive"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var payload struct {
		Answers []dto.QuizAnswerSubmission `json:"answers"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quiz answers"})
		return
	}

	questions, err := services.GetQuestions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	categoryQuestions := make([]models.Question, 0)
	for _, question := range questions {
		if question.CategoryId == categoryID {
			categoryQuestions = append(categoryQuestions, question)
		}
	}

	result, err := calculateQuizResult(categoryID, categoryQuestions, payload.Answers)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func calculateQuizResult(categoryID uint, questions []models.Question, submissions []dto.QuizAnswerSubmission) (dto.QuizResultDTO, error) {
	questionsByID := make(map[uint]models.Question, len(questions))
	for _, question := range questions {
		questionsByID[question.ID] = question
	}

	selectedByQuestion := make(map[uint]map[uint]struct{}, len(submissions))
	for _, submission := range submissions {
		question, ok := questionsByID[submission.QuestionID]
		if !ok {
			return dto.QuizResultDTO{}, fmt.Errorf("answer references a question outside this category")
		}
		if _, exists := selectedByQuestion[submission.QuestionID]; exists {
			return dto.QuizResultDTO{}, fmt.Errorf("each question can only be submitted once")
		}

		selected := make(map[uint]struct{}, len(submission.AnswerIDs))
		for _, answerID := range submission.AnswerIDs {
			if _, exists := selected[answerID]; exists {
				return dto.QuizResultDTO{}, fmt.Errorf("an answer cannot be selected more than once")
			}
			matchesQuestion := false
			for _, answer := range question.Answers {
				if answer.ID == answerID {
					matchesQuestion = true
					break
				}
			}
			if !matchesQuestion {
				return dto.QuizResultDTO{}, fmt.Errorf("answer does not belong to its question")
			}
			selected[answerID] = struct{}{}
		}
		selectedByQuestion[submission.QuestionID] = selected
	}

	correctCount := 0
	for _, question := range questions {
		selected := selectedByQuestion[question.ID]
		correctAnswers := 0
		allCorrectAnswersSelected := true
		for _, answer := range question.Answers {
			if answer.IsCorrect {
				correctAnswers++
				if _, ok := selected[answer.ID]; !ok {
					allCorrectAnswersSelected = false
				}
			}
		}
		if allCorrectAnswersSelected && len(selected) == correctAnswers {
			correctCount++
		}
	}

	percentage := 0
	if len(questions) > 0 {
		percentage = correctCount * 100 / len(questions)
	}
	return dto.QuizResultDTO{
		CategoryID:   categoryID,
		CorrectCount: correctCount,
		TotalCount:   len(questions),
		Percentage:   percentage,
	}, nil
}
