package model

import (
	"database/sql"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	DEVICE_INSERT_QUERY = "INSERT INTO devices (uuid, password, created_at) VALUES (?, ?, ?)"
	DEVICE_DELETE_QUERY = "DELETE FROM devices WHERE uuid = '?'"
)

type Device struct {
	UUID       uuid.UUID `json:"uuid"`
	DeviceId   int       `json:"deviceId"`
	Created_at time.Time `json:"createdAt"`
}

func MakeDevice(db *sql.DB, password []byte) (*Device, error) {

	// create UUID and hash password
	createdAt := time.Now()
	uid := uuid.Must(uuid.NewV4())
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// insert into database
	result, err := db.Exec(DEVICE_INSERT_QUERY, uid, hashedPassword, createdAt)
	if err != nil {
		return nil, err
	}
	deviceId, err := result.LastInsertId()

	if err != nil {
		// get device_id manually
		var id int
		err = db.QueryRow("SELECT device_id FROM devices ORDER BY created_at DESC LIMIT 1").Scan(&id)
		if err != nil {
			return nil, err
		}

		deviceId = int64(id)
	}

	return &Device{uid, int(deviceId), createdAt}, nil
}

func AuthDevice(db *sql.DB, uid, password string) (string, AuthToken, error) {

	return "", AuthToken{}, nil
}

func DeleteDeviceByUUID(db *sql.DB, uid string) error {

	_, err := db.Exec(DEVICE_DELETE_QUERY, uid)

	return err
}
