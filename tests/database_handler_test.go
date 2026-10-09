package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github/rizkypedia/quiz-master/internal/handlers"
	"github/rizkypedia/quiz-master/internal/security"
)

func TestCreateQuestionHandlerWithMockedDatabase(t *testing.T) {
	mock := mockDatabase(t)
	now := time.Now()
	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT \* FROM "categories" WHERE "categories"\."id" = \$1 AND "categories"\."deleted_at" IS NULL ORDER BY "categories"\."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name", "active"}).
			AddRow(1, now, now, nil, "Geography", true))
	mock.ExpectQuery(`SELECT \* FROM "quiz_types" WHERE "quiz_types"\."id" = \$1 AND "quiz_types"\."deleted_at" IS NULL ORDER BY "quiz_types"\."id" LIMIT \$2`).
		WithArgs(2, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"}).
			AddRow(2, now, now, nil, "SingleChoice"))
	mock.ExpectQuery(`INSERT INTO "questions".*RETURNING "id"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "Capital?", 1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))
	mock.ExpectQuery(`INSERT INTO "answers".*RETURNING "id"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "Paris", 10, true).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(20))
	mock.ExpectCommit()

	router := gin.New()
	router.POST("/question", handlers.CreateQuestion)
	request := httptest.NewRequest(http.MethodPost, "/question", strings.NewReader(
		`{"question":"Capital?","categoryId":1,"quizTypeId":2,"answers":[{"answer":"Paris","isCorrect":true}]}`,
	))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusCreated, response.Body)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}

func TestCreateCategoryHandlerWithMockedDatabase(t *testing.T) {
	mock := mockDatabase(t)
	mock.ExpectQuery(`INSERT INTO "categories".*RETURNING "id"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))

	router := gin.New()
	router.POST("/category", handlers.CreateCategory)
	request := httptest.NewRequest(http.MethodPost, "/category", strings.NewReader(`{"name":"Geography"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusCreated, response.Body)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}

func TestGetCategoriesReturnsOnlyActiveCategories(t *testing.T) {
	mock := mockDatabase(t)
	now := time.Now()
	mock.ExpectQuery(`SELECT \* FROM "categories" WHERE categories\.active = \$1 AND "categories"\."deleted_at" IS NULL ORDER BY categories\.id ASC`).
		WithArgs(true).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name", "active"}).
			AddRow(3, now, now, nil, "Geography", true))

	router := gin.New()
	router.GET("/categories", handlers.GetCategories)
	request := httptest.NewRequest(http.MethodGet, "/categories", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusOK, response.Body)
	}
	var body struct {
		Items []struct {
			ID     uint   `json:"id"`
			Name   string `json:"name"`
			Active bool   `json:"active"`
		} `json:"items"`
		Count int `json:"count"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Count != 1 || len(body.Items) != 1 || body.Items[0].Name != "Geography" || !body.Items[0].Active {
		t.Fatalf("unexpected active categories response: %+v", body)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}

func TestSetCategoryActiveHandler(t *testing.T) {
	mock := mockDatabase(t)
	mock.ExpectExec(`UPDATE "categories" SET "active"=\$1,"updated_at"=\$2 WHERE id = \$3 AND "categories"\."deleted_at" IS NULL`).
		WithArgs(false, sqlmock.AnyArg(), 3).
		WillReturnResult(sqlmock.NewResult(0, 1))

	router := gin.New()
	router.PUT("/category/:id/active", handlers.SetCategoryActive)
	request := httptest.NewRequest(http.MethodPut, "/category/3/active", strings.NewReader(`{"active":false}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusOK, response.Body)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}

func TestInactiveCategoryRejectsExistingQuizLink(t *testing.T) {
	mock := mockDatabase(t)
	t.Setenv("JWT_SECRET", "test-secret")
	token, _, err := security.IssueCategoryQuizLink(3)
	if err != nil {
		t.Fatalf("issue quiz link: %v", err)
	}
	mock.ExpectQuery(`SELECT \* FROM "categories" WHERE active = \$1 AND "categories"\."id" = \$2 AND "categories"\."deleted_at" IS NULL ORDER BY "categories"\."id" LIMIT \$3`).
		WithArgs(true, 3, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name", "active"}))

	router := gin.New()
	router.GET("/quiz", handlers.GetQuizByToken)
	request := httptest.NewRequest(http.MethodGet, "/quiz?token="+token, nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusUnauthorized, response.Body)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}

func TestGetQuizTypesHandlerWithMockedDatabase(t *testing.T) {
	mock := mockDatabase(t)
	now := time.Now()
	mock.ExpectQuery(`SELECT \* FROM "quiz_types" WHERE "quiz_types"\."deleted_at" IS NULL ORDER BY quiz_types\.id ASC`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"}).
			AddRow(1, now, now, nil, "MultipleChoice").
			AddRow(2, now, now, nil, "SingleChoice"))

	router := gin.New()
	router.GET("/quiz-types", handlers.GetQuizTypes)
	request := httptest.NewRequest(http.MethodGet, "/quiz-types", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusOK, response.Body)
	}
	var body struct {
		Items []struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		} `json:"items"`
		Count int `json:"count"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Count != 2 || len(body.Items) != 2 || body.Items[1].Name != "SingleChoice" {
		t.Fatalf("unexpected quiz type response: %+v", body)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}
