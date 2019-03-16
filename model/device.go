package model

import (
	"database/sql"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	INSERT_QUERY = "INSERT INTO devices (uuid, password, created_at) VALUES (?, ?, ?)"
	DELETE_QUERY = "DELETE FROM devices WHERE uuid = '?'"
)

type Device struct {
	UUID       uuid.UUID `json:"uuid"`
	DeviceId   int       `json:"deviceId"`
	Created_at time.Time `json:"createdAt"`
}

func MakeDevice(db *sql.DB, password string) (*Device, error) {

	// create UUID and hash password
	createdAt := time.Now()
	uid := uuid.Must(uuid.NewV4())
	hashedPassword := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	// insert into database
	rows, err := db.Query(INSERT_QUERY, uid, hashedPassword, createdAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// get device_id (is auto increment)
	var deviceId int
	err = db.QueryRow("SELECT device_id FROM devices ORDER BY created_at DESC LIMIT 1").Scan(&deviceId)
	if err != nil {
		return nil, err
	}

	return &Device{uid, deviceId, createdAt}, nil
}

func DeleteDeviceByUUID(db *sql.DB, uid string) error {

	err := device.DB.QueryRow(DELETE_QUERY, uid)

	return err
}
