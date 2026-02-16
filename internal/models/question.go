package models

import "gorm.io/gorm"

type Question struct {
	gorm.Model
	Question   string
	CategoryId uint
	Category   Category
	Answers    []Answers `gorm:"foreignKey:QuestionId;references:ID"`
}
