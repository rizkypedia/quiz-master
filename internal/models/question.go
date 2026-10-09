package models

import "gorm.io/gorm"

type Question struct {
	gorm.Model
	Question   string    `gorm:"not null"`
	CategoryId uint      `gorm:"not null"`
	Category   Category  `gorm:"foreignKey:CategoryId;references:ID"`
	QuizTypeID uint      `gorm:"not null"`
	QuizType   QuizType  `gorm:"foreignKey:QuizTypeID;references:ID"`
	Answers    []Answers `gorm:"foreignKey:QuestionId;references:ID"`
}
