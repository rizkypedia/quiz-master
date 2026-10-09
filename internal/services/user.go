package services

import (
	"errors"
	"fmt"
	"strings"

	"github/rizkypedia/quiz-master/internal/database"
	"github/rizkypedia/quiz-master/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidUser        = errors.New("invalid user input")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func CreateUser(email string, password string) error {
	email = strings.TrimSpace(email)
	if email == "" || password == "" {
		return fmt.Errorf("%w: email and password are required", ErrInvalidUser)
	}
	if len(email) > 255 {
		return fmt.Errorf("%w: email must not exceed 255 bytes", ErrInvalidUser)
	}
	if len([]byte(password)) > 72 {
		return fmt.Errorf("%w: password must not exceed 72 bytes", ErrInvalidUser)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user := models.User{
		Email:    email,
		Password: string(hashed),
	}
	return database.DB.Create(&user).Error
}

func Authenticate(email, password string) (*models.User, error) {
	email = strings.TrimSpace(email)
	if email == "" || password == "" {
		return nil, ErrInvalidCredentials
	}

	var user models.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, ErrInvalidCredentials
	}

	return &user, nil
}
