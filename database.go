package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func setupDatabase() *sql.DB {

	db, err := sql.Open("mysql", "homestead:secret@(mysql)/homestead?parseTime=true")
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
