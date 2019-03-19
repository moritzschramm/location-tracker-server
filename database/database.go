package database

import (
	"database/sql"
	"log"

	//_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

func SetupDatabase() *sql.DB {

	db, err := openWithSQLite3Driver()
	if err != nil {
		log.Fatal("Error opening database: ", err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err.Error())
	}

	return db
}

func openWithSQLite3Driver() (*sql.DB, error) {

	return sql.Open("sqlite3", "vault.db")
}

func openWithMySQLDriver() (*sql.DB, error) {	// if in use, uncomment driver import

	return sql.Open("mysql", "vault:secret@(mysql)/vault?parseTime=true")
}