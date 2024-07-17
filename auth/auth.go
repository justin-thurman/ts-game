package auth

import (
	"context"
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

type accountCreator interface {
	CreateAccount(ctx context.Context, username, passwordHash, salt string) (int32, error)
}

// CreateAccount creates a new player account.
func CreateAccount(ctx context.Context, creator accountCreator, username, password string) (int32, error) {
	salt, err := GenerateSalt()
	if err != nil {
		return -1, err
	}
	hashedPassword, err := HashPassword(password, salt)
	if err != nil {
		return -1, err
	}
	accountId, err := creator.CreateAccount(ctx, username, hashedPassword, salt)
	if err != nil {
		return -1, err
	}
	return accountId, nil
}
