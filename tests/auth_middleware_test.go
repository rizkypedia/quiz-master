package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github/rizkypedia/quiz-master/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TestAuthRequiredRejectsInvalidAuthorization(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	tests := []struct {
		name   string
		header string
	}{
		{name: "missing header"},
		{name: "invalid scheme", header: "Basic token"},
		{name: "invalid token", header: "Bearer invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/secure", middleware.AuthRequired(), func(c *gin.Context) {
				c.Status(http.StatusNoContent)
			})
			request := httptest.NewRequest(http.MethodGet, "/secure", nil)
			if tt.header != "" {
				request.Header.Set("Authorization", tt.header)
			}
			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)

			if response.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusUnauthorized, response.Body)
			}
		})
	}
}

func TestAuthRequiredAcceptsSignedBearerToken(t *testing.T) {
	const secret = "test-secret"
	t.Setenv("JWT_SECRET", secret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": 1})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign test token: %v", err)
	}

	router := gin.New()
	router.GET("/secure", middleware.AuthRequired(), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
	request := httptest.NewRequest(http.MethodGet, "/secure", nil)
	request.Header.Set("Authorization", "Bearer "+tokenString)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusNoContent, response.Body)
	}
}

func TestAuthRequiredRejectsExpiredToken(t *testing.T) {
	const secret = "test-secret"
	t.Setenv("JWT_SECRET", secret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"exp": 1})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign test token: %v", err)
	}

	router := gin.New()
	router.GET("/secure", middleware.AuthRequired(), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
	request := httptest.NewRequest(http.MethodGet, "/secure", nil)
	request.Header.Set("Authorization", "Bearer "+tokenString)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusUnauthorized, response.Body)
	}
}
