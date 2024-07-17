package auth

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// GenerateSalt creates a random salt.
func GenerateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

// HashPassword hashes a password using the provided salt.
func HashPassword(password, salt string) (string, error) {
	saltedPassword := password + salt
	hash, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePassword compares a plaintext password against the hashed and salted password.
func ComparePassword(plaintextPassword, hashedPassword, salt string) (bool, error) {
	plaintextPlusSalt := plaintextPassword + salt
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plaintextPlusSalt))
	if err != nil {
		return false, err
	}
	return true, nil
}
