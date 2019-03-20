package database

import (
	"database/sql"
	"log"
	"io/ioutil"

	//_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

const (
	DB_FILE_SQLITE = "database/vault.db"
	DB_INIT_STMT = "database/init_database.sql"
)

func SetupDatabase() *sql.DB {

	// create database interface
	db, err := openWithSQLite3Driver()
	if err != nil {
		log.Fatal("Error opening database: ", err.Error())
	}

	// check connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err.Error())
	}

	// read init statement from file DB_INIT_STMT
	initStmt, err := ioutil.ReadFile(DB_INIT_STMT)
	if err != nil {
		log.Fatal("Error reading database init statement file: ", err.Error())
	}

	// init database tables
	_, err = db.Exec(string(initStmt))
	if err != nil {
		log.Fatal("Error executing init statement: ", err.Error())
	}

	return db
}

func openWithSQLite3Driver() (*sql.DB, error) {

	return sql.Open("sqlite3", DB_FILE_SQLITE)
}

func openWithMySQLDriver() (*sql.DB, error) {	// if in use, uncomment driver import

	return sql.Open("mysql", "vault:secret@(mysql)/vault?parseTime=true")
}