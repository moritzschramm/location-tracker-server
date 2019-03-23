package model

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/satori/go.uuid"
)

const (
	QUERY_DEVICE_BY_ID = "SELECT uuid, created_at FROM devices WHERE device_id == ?"
	INSERT_DEVICE      = "INSERT INTO devices (uuid, password, created_at) VALUES (?, ?, ?)"
	DELETE_DEVICE      = "DELETE FROM devices WHERE uuid = '?'"
)

type Device struct {
	UUID      uuid.UUID `json:"uuid"`
	DeviceId  int       `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
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
	result, err := db.Exec(INSERT_DEVICE, uid, hashedPassword, createdAt)
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

func GetDevice(db *sql.DB, deviceId int) (*Device, error) {

	var uidString string
	var createdAt time.Time
	err := db.QueryRow(QUERY_DEVICE_BY_ID, deviceId).Scan(&uidString, &createdAt)
	if err != nil {
		return nil, err
	}

	uid, _ := uuid.FromString(uidString)

	return &Device{
		UUID:      uid,
		DeviceId:  deviceId,
		CreatedAt: createdAt,
	}, nil
}

func DeleteDeviceByUUID(db *sql.DB, uid string) error {

	_, err := db.Exec(DELETE_DEVICE, uid)

	return err
}
