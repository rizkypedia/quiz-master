package models

import (
	"gorm.io/gorm"
)

type Answers struct {
	gorm.Model
	Answer     string `gorm:"not null"`
	QuestionId uint   `gorm:"not null"`
	IsCorrect  bool   `gorm:"not null;default:false"`
}
