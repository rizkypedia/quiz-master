package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github/rizkypedia/quiz-master/internal/handlers"

	"github.com/gin-gonic/gin"
)

func TestCreateQuestionRejectsMissingQuizType(t *testing.T) {
	router := gin.New()
	router.POST("/question", handlers.CreateQuestion)
	request := httptest.NewRequest(http.MethodPost, "/question", strings.NewReader(
		`{"question":"Question?","categoryId":1,"answers":[{"answer":"Answer","isCorrect":true}]}`,
	))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusBadRequest, response.Body)
	}
}

func TestCreateQuestionRejectsMalformedJSON(t *testing.T) {
	router := gin.New()
	router.POST("/question", handlers.CreateQuestion)
	request := httptest.NewRequest(http.MethodPost, "/question", strings.NewReader(`{"question":`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusBadRequest, response.Body)
	}
}

func TestCreateCategoryRequiresName(t *testing.T) {
	router := gin.New()
	router.POST("/category", handlers.CreateCategory)
	request := httptest.NewRequest(http.MethodPost, "/category", strings.NewReader(`{}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusBadRequest, response.Body)
	}
}

func TestCreateCategoryRejectsWhitespaceName(t *testing.T) {
	router := gin.New()
	router.POST("/category", handlers.CreateCategory)
	request := httptest.NewRequest(http.MethodPost, "/category", strings.NewReader(`{"name":"   "}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusBadRequest, response.Body)
	}
}

func TestRegisterAndLoginRejectInvalidRequests(t *testing.T) {
	tests := []struct {
		name    string
		handler gin.HandlerFunc
		body    string
	}{
		{name: "register malformed JSON", handler: handlers.Register, body: `{"email":`},
		{name: "register invalid email", handler: handlers.Register, body: `{"email":"not-an-email","password":"secret"}`},
		{name: "login missing password", handler: handlers.Login, body: `{"email":"user@example.com"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/auth", tt.handler)
			request := httptest.NewRequest(http.MethodPost, "/auth", strings.NewReader(tt.body))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)

			if response.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusBadRequest, response.Body)
			}
		})
	}
}
