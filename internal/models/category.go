package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name   string `gorm:"not null"`
	Active bool   `gorm:"not null;default:true"`
}
