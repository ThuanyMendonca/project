package utils

import "golang.org/x/crypto/bcrypt"

func GenerateHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func CompareHash(passwordHash, passwordString string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordString))
}
