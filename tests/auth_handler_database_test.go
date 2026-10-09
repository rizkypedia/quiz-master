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
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github/rizkypedia/quiz-master/internal/handlers"
)

func TestRegisterHandlerWithMockedDatabase(t *testing.T) {
	mock := mockDatabase(t)
	mock.ExpectQuery(`INSERT INTO "users".*RETURNING "id"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "user@example.com", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(5))

	router := gin.New()
	router.POST("/register", handlers.Register)
	request := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(
		`{"email":"user@example.com","password":"secure-password"}`,
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

func TestLoginHandlerWithMockedDatabase(t *testing.T) {
	const secret = "test-secret"
	t.Setenv("JWT_SECRET", secret)
	mock := mockDatabase(t)
	hash, err := bcrypt.GenerateFromPassword([]byte("secure-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash test password: %v", err)
	}
	now := time.Now()
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs("user@example.com", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "password"}).
			AddRow(5, now, now, nil, "user@example.com", string(hash)))

	router := gin.New()
	router.POST("/login", handlers.Login)
	request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(
		`{"email":"user@example.com","password":"secure-password"}`,
	))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusOK, response.Body)
	}
	var body struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	parsedToken, err := jwt.Parse(body.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !parsedToken.Valid {
		t.Fatalf("returned token is invalid: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}
