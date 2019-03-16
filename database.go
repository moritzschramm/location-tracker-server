package main

import (
	"database/sql"

	// _ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

const (
	DRIVER = "sqlite3"	// "mysql"
	MYSQL_DSN = "vault:secret@(mysql)/vault?parseTime=true"
	SQLITE_DSN = "vault.db"
)

func setupDatabase() *sql.DB {

	db, err := sql.Open(DRIVER, SQLITE_DSN)
	if err != nil {
		log.Fatal("Error opening database: ", err.Error())
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err.Error())
		panic(err)
	}

	return db
}
