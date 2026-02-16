package models

import "gorm.io/gorm"

type QuizType struct {
	gorm.Model
	Name string
}
