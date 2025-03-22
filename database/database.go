package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDBConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	return db, nil
}
