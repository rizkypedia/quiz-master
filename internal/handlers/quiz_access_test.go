package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github/rizkypedia/quiz-master/internal/dto"
	"github/rizkypedia/quiz-master/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestCalculateQuizResult(t *testing.T) {
	questions := []models.Question{
		{
			Model:      gorm.Model{ID: 10},
			CategoryId: 1,
			Answers: []models.Answers{
				{Model: gorm.Model{ID: 101}, IsCorrect: true},
				{Model: gorm.Model{ID: 102}},
			},
		},
		{
			Model:      gorm.Model{ID: 20},
			CategoryId: 1,
			Answers: []models.Answers{
				{Model: gorm.Model{ID: 201}, IsCorrect: true},
				{Model: gorm.Model{ID: 202}, IsCorrect: true},
				{Model: gorm.Model{ID: 203}},
			},
		},
	}

	result, err := calculateQuizResult(1, questions, []dto.QuizAnswerSubmission{
		{QuestionID: 10, AnswerIDs: []uint{101}},
		{QuestionID: 20, AnswerIDs: []uint{201, 202}},
	})
	if err != nil {
		t.Fatalf("calculateQuizResult() error = %v", err)
	}
	if result.CategoryID != 1 || result.CorrectCount != 2 || result.TotalCount != 2 || result.Percentage != 100 {
		t.Fatalf("unexpected quiz result: %+v", result)
	}
}

func TestCalculateQuizResultRequiresExactMultipleChoiceSet(t *testing.T) {
	questions := []models.Question{{
		Model: gorm.Model{ID: 20},
		Answers: []models.Answers{
			{Model: gorm.Model{ID: 201}, IsCorrect: true},
			{Model: gorm.Model{ID: 202}, IsCorrect: true},
			{Model: gorm.Model{ID: 203}},
		},
	}}

	result, err := calculateQuizResult(1, questions, []dto.QuizAnswerSubmission{
		{QuestionID: 20, AnswerIDs: []uint{201}},
	})
	if err != nil {
		t.Fatalf("calculateQuizResult() error = %v", err)
	}
	if result.CorrectCount != 0 || result.Percentage != 0 {
		t.Fatalf("partial multiple-choice selection should be incorrect: %+v", result)
	}
}

func TestCalculateQuizResultRejectsAnswersOutsideQuiz(t *testing.T) {
	questions := []models.Question{{
		Model: gorm.Model{ID: 10},
		Answers: []models.Answers{
			{Model: gorm.Model{ID: 101}, IsCorrect: true},
		},
	}}

	_, err := calculateQuizResult(1, questions, []dto.QuizAnswerSubmission{
		{QuestionID: 99, AnswerIDs: []uint{101}},
	})
	if err == nil {
		t.Fatal("expected an error for an answer to a question outside the quiz")
	}

	_, err = calculateQuizResult(1, questions, []dto.QuizAnswerSubmission{
		{QuestionID: 10, AnswerIDs: []uint{999}},
	})
	if err == nil {
		t.Fatal("expected an error for an answer that does not belong to the question")
	}
}

func TestCalculateQuizResultCountsUnansweredQuestionsAsIncorrect(t *testing.T) {
	questions := []models.Question{
		{
			Model: gorm.Model{ID: 10},
			Answers: []models.Answers{
				{Model: gorm.Model{ID: 101}, IsCorrect: true},
			},
		},
		{
			Model: gorm.Model{ID: 20},
			Answers: []models.Answers{
				{Model: gorm.Model{ID: 201}, IsCorrect: true},
			},
		},
	}

	result, err := calculateQuizResult(1, questions, []dto.QuizAnswerSubmission{
		{QuestionID: 10, AnswerIDs: []uint{101}},
	})
	if err != nil {
		t.Fatalf("calculateQuizResult() error = %v", err)
	}
	if result.CorrectCount != 1 || result.TotalCount != 2 || result.Percentage != 50 {
		t.Fatalf("unexpected quiz result with an unanswered question: %+v", result)
	}
}

func TestQuizEndpointsRequireSecureLink(t *testing.T) {
	router := gin.New()
	router.GET("/quiz", GetQuizByToken)
	router.POST("/quiz/result", SubmitQuizResult)

	for _, endpoint := range []struct {
		method string
		path   string
	}{
		{method: http.MethodGet, path: "/quiz"},
		{method: http.MethodPost, path: "/quiz/result"},
	} {
		request := httptest.NewRequest(endpoint.method, endpoint.path, strings.NewReader(`{}`))
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		if response.Code != http.StatusUnauthorized {
			t.Errorf("%s %s status = %d, want %d", endpoint.method, endpoint.path, response.Code, http.StatusUnauthorized)
		}
	}
}

func TestAnswerResponseDoesNotExposeCorrectness(t *testing.T) {
	response, err := json.Marshal(dto.AnswerDTO{ID: 1, Answer: "Correct answer", QuestionId: 2})
	if err != nil {
		t.Fatalf("marshal answer DTO: %v", err)
	}
	if strings.Contains(string(response), "isCorrect") {
		t.Fatalf("answer response exposed correctness: %s", response)
	}
}
