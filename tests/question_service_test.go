package tests

import (
	"errors"
	"testing"

	"github/rizkypedia/quiz-master/internal/models"
	"github/rizkypedia/quiz-master/internal/services"
)

func TestCreateQuestionRejectsInvalidInputBeforeDatabaseAccess(t *testing.T) {
	tests := []struct {
		name     string
		question *models.Question
	}{
		{name: "nil question"},
		{name: "blank question text", question: &models.Question{Question: "  ", CategoryId: 1, QuizTypeID: 1}},
		{name: "missing category", question: &models.Question{Question: "Question?", QuizTypeID: 1}},
		{name: "missing quiz type", question: &models.Question{Question: "Question?", CategoryId: 1}},
		{name: "missing answers", question: &models.Question{Question: "Question?", CategoryId: 1, QuizTypeID: 1}},
		{
			name: "blank answer",
			question: &models.Question{
				Question: "Question?", CategoryId: 1, QuizTypeID: 1,
				Answers: []models.Answers{{Answer: "  "}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := services.CreateQuestion(tt.question)
			if !errors.Is(err, services.ErrInvalidQuestion) {
				t.Fatalf("CreateQuestion() error = %v, want ErrInvalidQuestion", err)
			}
		})
	}
}
