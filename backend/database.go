package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	dbName string
	db     *sql.DB
}

func (db *Database) init() error {
	var err error
	db.db, err = sql.Open("sqlite3", db.dbName)

	return err
}

func (db *Database) close() {
	if db.db != nil {
		db.db.Close()
	}
}

// func (db *Database) (query string, ) {
// 	db.db.Query()
// }
