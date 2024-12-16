package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// HashedPassword returns a bcrypt hash password
func HashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate hashed password:%s", err)
	}
	return string(hashedPassword), nil
}

func checkPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
