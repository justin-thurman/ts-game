package auth

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"ts-game/db/queries"
)

type loginSignUpHandler interface {
	CreateAccount(ctx context.Context, username, passwordHash, salt string) (int32, error)
	AccountExists(ctx context.Context, username string) (bool, error)
	GetAccount(ctx context.Context, username string) (queries.GetAccountRow, error)
}

// LoginOrSignUp handles the flow for a user logging into an existing account or signing up for a new one.
func LoginOrSignUp(ctx context.Context, r io.Reader, w io.Writer, handler loginSignUpHandler) int32 {
	scanner := bufio.NewScanner(r)
	var accountId int32
	for scanner.Scan() {
		accountName := scanner.Text()
		accountExists, err := AccountExists(ctx, handler, accountName)
		if err != nil {
			fmt.Fprintln(w, "Error searching for account. Please try again.")
			continue
		}
		if accountExists {
			fmt.Fprintf(w, "Logging into account %s. Enter password.\n", accountName)
			scanner.Scan()
			password := scanner.Text()
			accountId, err = Login(ctx, handler, accountName, password)
			if err != nil {
				if err.Error() == "incorrect password" {
					fmt.Fprintln(w, "Incorrect password.")
				} else {
					fmt.Fprintln(w, "Error logging in. Please try again.")
					slog.Error("Error during login", "err", err, "accountName", accountName)
				}
				fmt.Fprintln(w, "Welcome! Enter your account name to login to an existing account or create a new one.")
				continue
			}
			slog.Debug("Login to account", "accountId", accountId)
			fmt.Fprintf(w, "Welcome back, %s!", accountName)
			break
		} else {
			fmt.Fprintf(w, "Creating account with name %s. Continue? 'yes' or 'no'\n", accountName)
			scanner.Scan()
			answer := scanner.Text()
			answer = strings.ToLower(answer)
			if answer != "yes" {
				fmt.Fprintln(w, "Welcome! Enter your account name to login to an existing account or create a new one.")
				continue
			}
			password := ""
			password2 := "not matching"
			for password != password2 {
				fmt.Fprintln(w, "Please enter your password.")
				scanner.Scan()
				password = scanner.Text()
				fmt.Fprintln(w, "Please enter your password again.")
				scanner.Scan()
				password2 = scanner.Text()
				if password != password2 {
					fmt.Fprintln(w, "Passwords do not match")
				}
			}
			accountId, err = CreateAccount(ctx, handler, accountName, password)
			if err != nil {
				fmt.Fprintln(w, "Error creating account. Please try again.")
				slog.Error("Error during account creation", "err", err)
			}
			slog.Debug("Account created", "accountId", accountId)
			fmt.Fprintln(w, "Account created successfully!")
			break
		}
	}
	return accountId
}
