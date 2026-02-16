package models

import (
	"gorm.io/gorm"
)

type Answers struct {
	gorm.Model
	Answer     string
	QuestionId uint
	IsCorrect  bool
}
