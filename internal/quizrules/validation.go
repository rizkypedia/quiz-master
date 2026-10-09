package quizrules

import (
	"errors"
	"fmt"
	"strings"

	"github/rizkypedia/quiz-master/internal/models"
)

var ErrInvalidQuestion = errors.New("invalid question")

func ValidateAnswers(answers []models.Answers) error {
	if len(answers) == 0 {
		return fmt.Errorf("%w: at least one answer is required", ErrInvalidQuestion)
	}
	for i := range answers {
		answers[i].Answer = strings.TrimSpace(answers[i].Answer)
		if answers[i].Answer == "" {
			return fmt.Errorf("%w: answer text is required", ErrInvalidQuestion)
		}
	}
	return nil
}

func ValidateCorrectAnswerCount(quizType string, answers []models.Answers) error {
	correct := 0
	for _, answer := range answers {
		if answer.IsCorrect {
			correct++
		}
	}

	quizType = strings.TrimSpace(quizType)
	singleChoice := strings.EqualFold(quizType, "SingleChoice") ||
		strings.EqualFold(quizType, "SingleChoise")
	multipleChoice := strings.EqualFold(quizType, "MultipleChoice")

	switch {
	case singleChoice && correct != 1:
		return fmt.Errorf("%w: single-choice questions must have exactly one correct answer", ErrInvalidQuestion)
	case multipleChoice && correct == 0:
		return fmt.Errorf("%w: multiple-choice questions must have at least one correct answer", ErrInvalidQuestion)
	case !singleChoice && !multipleChoice:
		return fmt.Errorf("%w: unsupported quiz type %q", ErrInvalidQuestion, quizType)
	default:
		return nil
	}
}
