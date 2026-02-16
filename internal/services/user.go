package services

import (
	"github/rizkypedia/quiz-master/internal/database"
	"github/rizkypedia/quiz-master/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(email string, password string) error {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	user := models.User{
		Email:    email,
		Password: string(hashed),
	}
	return database.DB.Create(&user).Error
}

func Authenticate(email, password string) (*models.User, bool) {
	var user models.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, false
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, false
	}

	return &user, true
}
