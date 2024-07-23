package main

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Account struct {
	ID    int
	Token string
}

var (
	ErrWrongPassword = errors.New("wrong password")
	ErrSqlFailed     = errors.New("sql failed")
	ErrNoAccount     = errors.New("no account found with that email")
)

func accountFromCredentials(email, password string) (*Account, error) {
	hashedSaltedPassword := password // TODO
	row := db.QueryRow(
		"SELECT id, hash, token FROM users WHERE email = ?",
		email,
	)

	var id int
	var hash string
	var token string

	if err := row.Scan(&id, &hash, &token); err == sql.ErrNoRows {
		return nil, ErrNoAccount
	} else if err != nil {
		return nil, ErrSqlFailed
	}

	if hashedSaltedPassword == hash {
		return &Account{ID: id, Token: token}, nil
	}

	return nil, ErrWrongPassword
}

// DOES NOT CHECK VALIDATION
func addTierToAccount(accIDString, tierString string) (string, error) {
	tx, err := db.Begin()

	if err != nil {
		return "", err
	}

	timeRN := time.Now().UnixMilli()
	generatedURL := "change-this-pls" // TODO: should be obvious

	_, err = tx.Exec(
		`INSERT INTO tiers (id, accID, tier, url) VALUES (?, ?, ?, ?)`, timeRN, accIDString, tierString, generatedURL,
	)

	if err != nil {
		tx.Rollback()
		return "", err
	}

	tx.Commit()

	return generatedURL, nil
}

// Return whether the token corresponds to the given [idString] or not
func authenticate(idString, token string) bool {
	row := db.QueryRow("SELECT token FROM users WHERE id = ?", idString)
	var fetchedToken string

	if err := row.Scan(&fetchedToken); err == sql.ErrNoRows {
		return false
	} else if token == fetchedToken {
		return true
	}

	return false
}

// DOES NOT CHECK VALIDATION
func deleteFromAccount(tierID string) error {
	tx, err := db.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM tiers WHERE tierID = ?`, tierID)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
