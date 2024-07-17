package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"ts-game/db/queries"

	"golang.org/x/crypto/bcrypt"
)

func generateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

func hashPassword(password, salt string) (string, error) {
	saltedPassword := password + salt
	hash, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func comparePassword(plaintextPassword, hashedPassword, salt string) error {
	plaintextPlusSalt := plaintextPassword + salt
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plaintextPlusSalt))
}

type accountCreator interface {
	CreateAccount(ctx context.Context, username, passwordHash, salt string) (int32, error)
}

// CreateAccount creates a new player account.
func CreateAccount(ctx context.Context, creator accountCreator, username, password string) (int32, error) {
	salt, err := generateSalt()
	if err != nil {
		return -1, err
	}
	hashedPassword, err := hashPassword(password, salt)
	if err != nil {
		return -1, err
	}
	accountId, err := creator.CreateAccount(ctx, username, hashedPassword, salt)
	if err != nil {
		return -1, err
	}
	return accountId, nil
}

type accountExistsChecker interface {
	AccountExists(ctx context.Context, username string) (bool, error)
}

// AccountExists checks whether an account with the given username exists.
func AccountExists(ctx context.Context, accountExistsChecker accountExistsChecker, username string) (bool, error) {
	exists, err := accountExistsChecker.AccountExists(ctx, username)
	if err != nil {
		return false, err
	}
	return exists, nil
}

type accountGetter interface {
	GetAccount(ctx context.Context, username string) (queries.GetAccountRow, error)
}

// Login logs a user in to an existing account if the correct password is provided.
func Login(ctx context.Context, accountGetter accountGetter, username, password string) (int32, error) {
	account, err := accountGetter.GetAccount(ctx, username)
	if err != nil {
		return -1, err
	}
	err = comparePassword(password, account.PasswordHash, account.Salt)
	if err != nil {
		return -1, errors.New("incorrect password")
	}
	return account.ID, nil
}
