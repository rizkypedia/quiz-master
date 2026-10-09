package tests

import (
	"errors"
	"testing"

	"github/rizkypedia/quiz-master/internal/models"
	"github/rizkypedia/quiz-master/internal/quizrules"
)

func TestValidateCorrectAnswerCount(t *testing.T) {
	tests := []struct {
		name     string
		quizType string
		answers  []models.Answers
		wantErr  bool
	}{
		{
			name:     "single choice with one correct answer",
			quizType: "SingleChoice",
			answers:  []models.Answers{{IsCorrect: true}, {IsCorrect: false}},
		},
		{
			name:     "legacy misspelled single choice type",
			quizType: "SingleChoise",
			answers:  []models.Answers{{IsCorrect: true}, {IsCorrect: false}},
		},
		{
			name:     "single choice with multiple correct answers",
			quizType: "SingleChoice",
			answers:  []models.Answers{{IsCorrect: true}, {IsCorrect: true}},
			wantErr:  true,
		},
		{
			name:     "single choice without a correct answer",
			quizType: "SingleChoice",
			answers:  []models.Answers{{IsCorrect: false}},
			wantErr:  true,
		},
		{
			name:     "multiple choice with multiple correct answers",
			quizType: "MultipleChoice",
			answers:  []models.Answers{{IsCorrect: true}, {IsCorrect: true}},
		},
		{
			name:     "multiple choice requires a correct answer",
			quizType: "MultipleChoice",
			answers:  []models.Answers{{IsCorrect: false}},
			wantErr:  true,
		},
		{
			name:     "unsupported quiz type",
			quizType: "Essay",
			answers:  []models.Answers{{IsCorrect: true}},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := quizrules.ValidateCorrectAnswerCount(tt.quizType, tt.answers)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateCorrectAnswerCount() error = %v, wantErr %t", err, tt.wantErr)
			}
			if tt.wantErr && !errors.Is(err, quizrules.ErrInvalidQuestion) {
				t.Errorf("error %v does not wrap ErrInvalidQuestion", err)
			}
		})
	}
}

func TestValidateAnswers(t *testing.T) {
	tests := []struct {
		name    string
		answers []models.Answers
		wantErr bool
	}{
		{name: "empty answer list", wantErr: true},
		{name: "blank answer", answers: []models.Answers{{Answer: "  "}}, wantErr: true},
		{name: "valid answer", answers: []models.Answers{{Answer: "  yes  "}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := quizrules.ValidateAnswers(tt.answers)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateAnswers() error = %v, wantErr %t", err, tt.wantErr)
			}
			if !tt.wantErr && tt.answers[0].Answer != "yes" {
				t.Errorf("answer text was not trimmed: %q", tt.answers[0].Answer)
			}
		})
	}
}
