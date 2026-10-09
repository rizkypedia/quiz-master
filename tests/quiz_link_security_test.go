package tests

import (
	"testing"
	"time"

	"github/rizkypedia/quiz-master/internal/security"
)

func TestIssueAndValidateCategoryQuizLink(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	tokenString, expiresAt, err := security.IssueCategoryQuizLink(42)
	if err != nil {
		t.Fatalf("issue quiz link: %v", err)
	}

	if expiresAt.Sub(time.Now().UTC()) < 11*time.Hour {
		t.Fatalf("token should last about 12 hours, got expiry %s", expiresAt)
	}

	categoryID, validatedAt, err := security.ValidateCategoryQuizLink(tokenString)
	if err != nil {
		t.Fatalf("validate quiz link: %v", err)
	}

	if categoryID != 42 {
		t.Fatalf("expected category id 42, got %d", categoryID)
	}

	if validatedAt.IsZero() {
		t.Fatal("expected a valid expiry timestamp")
	}
}
