package main

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type Account struct {
	ID    int
	Token string
	// Email   string
	// Session string
}

var (
	ErrWrongPassword = errors.New("wrong password")
	ErrSqlFailed     = errors.New("sql failed")
	ErrNoAccount     = errors.New("no account found with that email")
)

func accountFromCredentials(email, password string) (*Account, error) {
	hashedSaltedPassword := password // TODO
	row := db.db.QueryRow("SELECT id, hash, token FROM users WHERE email = ?", email)

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
