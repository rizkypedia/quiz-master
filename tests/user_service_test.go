package tests

import (
	"errors"
	"strings"
	"testing"

	"github/rizkypedia/quiz-master/internal/services"
)

func TestCreateUserRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
	}{
		{name: "missing email", password: "password"},
		{name: "missing password", email: "user@example.com"},
		{name: "email exceeds column length", email: strings.Repeat("a", 247) + "@test.com", password: "password"},
		{name: "password exceeds bcrypt limit", email: "user@example.com", password: strings.Repeat("p", 73)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := services.CreateUser(tt.email, tt.password)
			if !errors.Is(err, services.ErrInvalidUser) {
				t.Fatalf("CreateUser() error = %v, want ErrInvalidUser", err)
			}
		})
	}
}

func TestAuthenticateRejectsMissingCredentials(t *testing.T) {
	_, err := services.Authenticate("", "password")
	if !errors.Is(err, services.ErrInvalidCredentials) {
		t.Fatalf("Authenticate() error = %v, want ErrInvalidCredentials", err)
	}
}
