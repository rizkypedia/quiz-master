package handlers

import (
	"github/rizkypedia/quiz-master/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type QuizTypeDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CategoryDTO struct {
	ID       uint        `json:"id"`
	Name     string      `json:"name"`
	QuizType QuizTypeDTO `json:"quizType"`
}

type AnswerDTO struct {
	ID         uint   `json:"id"`
	Answer     string `json:"answer"`
	QuestionId uint   `json:"questionId"`
	IsCorrect  bool   `json:"isCorrect"`
}

type QuestionDTO struct {
	ID         uint        `json:"id"`
	Question   string      `json:"question"`
	CategoryId uint        `json:"categoryId"`
	Category   CategoryDTO `json:"category"`
	Answers    []AnswerDTO `json:"answers"`
}

func GetQuestions(c *gin.Context) {
	qs, err := services.GetQuestions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	dtos := make([]QuestionDTO, len(qs))
	for i, q := range qs {
		answers := make([]AnswerDTO, len(q.Answers))
		for j, a := range q.Answers {
			answers[j] = AnswerDTO{
				ID: a.ID, Answer: a.Answer, QuestionId: a.QuestionId, IsCorrect: a.IsCorrect,
			}
		}
		dtos[i] = QuestionDTO{
			ID: q.ID, Question: q.Question, CategoryId: q.CategoryId,
			Category: CategoryDTO{
				ID: q.Category.ID, Name: q.Category.Name,
				QuizType: QuizTypeDTO{
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

func AddQuestion(c *gin.Context) {
	var body struct {
		Question string `json:"question"`
	}
	c.BindJSON(&body)

	if err := services.CreateQuestion(body.Question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "created"})
}
