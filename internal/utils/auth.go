package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("can't hash password: %w", err)
	}
	return string(passwordHash), nil
}

func CheckPassword(passwordHash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}

func GenerateToken() string {
	randomToken := make([]byte, 32)
	rand.Read(randomToken)
	return base64.URLEncoding.EncodeToString(randomToken)
}
