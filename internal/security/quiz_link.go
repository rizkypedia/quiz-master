package security

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const quizLinkTTL = 12 * time.Hour

func IssueCategoryQuizLink(categoryID uint) (string, time.Time, error) {
	if categoryID == 0 {
		return "", time.Time{}, fmt.Errorf("category id is invalid")
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", time.Time{}, fmt.Errorf("JWT_SECRET not set")
	}

	now := time.Now().UTC()
	expiresAt := now.Add(quizLinkTTL)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"category_id": categoryID,
		"exp":         expiresAt.Unix(),
		"iat":         now.Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign quiz token: %w", err)
	}

	return tokenString, expiresAt, nil
}

func ValidateCategoryQuizLink(tokenString string) (uint, time.Time, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return 0, time.Time{}, fmt.Errorf("JWT_SECRET not set")
	}

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil || !parsedToken.Valid {
		return 0, time.Time{}, fmt.Errorf("invalid quiz token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, time.Time{}, fmt.Errorf("unexpected token claims")
	}

	quizIDValue, ok := claims["category_id"]
	if !ok {
		return 0, time.Time{}, fmt.Errorf("category id missing from token")
	}

	categoryIDFloat, ok := quizIDValue.(float64)
	if !ok {
		return 0, time.Time{}, fmt.Errorf("category id has an invalid format")
	}

	categoryID := uint(categoryIDFloat)
	if categoryID == 0 {
		return 0, time.Time{}, fmt.Errorf("category id is invalid")
	}

	expValue, ok := claims["exp"].(float64)
	if !ok {
		return 0, time.Time{}, fmt.Errorf("token expiry is missing")
	}

	expiresAt := time.Unix(int64(expValue), 0).UTC()
	if time.Now().UTC().After(expiresAt) {
		return 0, expiresAt, fmt.Errorf("quiz link expired")
	}

	return categoryID, expiresAt, nil
}
