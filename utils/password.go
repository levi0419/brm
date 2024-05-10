package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashText returns the bcrypt hash of any text
func HashText(text string) (string, error) {
	hashedText, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedText), nil
}

// CheckHashedText checks if the provided plain text is correct or not
func CheckHashedText(plainText string, hashedText string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedText), []byte(plainText))
}
