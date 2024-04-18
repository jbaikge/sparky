package password

import (
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

func TemporaryPassword(length int) string {
	letters := []byte(`ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789`)
	pw := make([]byte, length)
	for i := range pw {
		pw[i] = letters[rand.Intn(len(letters))]
	}
	return string(pw)
}

func Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

func Validate(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
