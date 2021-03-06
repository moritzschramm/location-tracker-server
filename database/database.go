package database

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"time"

	"github.com/moritzschramm/location-tracker-server/config"

	//_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

const (
	DB_INIT_STMT  = "database/init_database.sql"
	DB_DSN_SQLITE = "database/vault.db"
	DB_DSN_MYSQL  = "vault:secret@(mysql)/vault?parseTime=true"

	QUERY_ADMIN  = "SELECT device_id FROM devices WHERE uuid = ?"
	INSERT_ADMIN = "INSERT INTO devices (device_id, uuid, password, created_at) VALUES (?, ?, ?, ?)"
)

// open and check database connection and init tables
func Setup(config config.Config) *sql.DB {

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

	// init database tables
	err = initTables(db, config)
	if err != nil {
		log.Fatal("Error executing database init statement: ", err.Error())
	}

	return db
}

func openWithSQLite3Driver() (*sql.DB, error) {

	return sql.Open("sqlite3", DB_DSN_SQLITE)
}

func openWithMySQLDriver() (*sql.DB, error) { // if in use, uncomment driver import

	return sql.Open("mysql", DB_DSN_MYSQL)
}

// execute init database sql statement (to create tables)
// and insert admin user as specified in config file
func initTables(db *sql.DB, config config.Config) error {

	// read init statement from file DB_INIT_STMT
	initStmt, err := ioutil.ReadFile(DB_INIT_STMT)
	if err != nil {
		log.Fatal("Error reading database init statement file: ", err.Error())
	}

	// create tables (if not already present)
	_, err = db.Exec(string(initStmt))
	if err != nil {
		return err
	}

	// query for admin user, if not existing, insert new user ("admin device")
	var deviceId int
	err = db.QueryRow(QUERY_ADMIN, config.AdminUUID).Scan(&deviceId)
	if err != nil {

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(config.AdminPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		deviceId = 1
		_, err = db.Exec(INSERT_ADMIN, deviceId, config.AdminUUID, hashedPassword, time.Now())
		if err != nil {
			return err
		}
	}

	return nil
}
