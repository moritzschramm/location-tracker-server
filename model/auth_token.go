package model

import (
	"golang.org/x/crypto/bcrypt"
	"database/sql"
	"time"

	"github.com/satori/go.uuid"
)

const (
	QUERY_TOKEN = "SELECT device_id FROM tokens WHERE token = '?' AND expires_at <= ?"
	QUERY_DEVICE_BY_ID = "SELECT uuid, created_at FROM devices WHERE device_id == ?"
	QUERY_DEVICE_BY_UUID = "SELECT device_id, password FROM devices WHERE uuid == ?"
)

type AuthToken struct {
	Id int 		 		`json:"-"`
	DeviceId int 		`json:"-"`
	Token string 		`json:"token"`
	CreatedAt time.Time `json:"-"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func AuthDevice(db *sql.DB, uid, password string) (*AuthToken, error) {

	// get password of device with uuid from database
	var deviceId int
	var hash string
	err := db.QueryRow(QUERY_DEVICE_BY_UUID, uid).Scan(&deviceId, &hash)
	if err != nil {
		return nil, err
	}

	// check if password hashes match
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return nil, err
	}

	// create new random token (256 bit long)
	// TODO
	token := "hello"

	// insert new token in database
	var id int

	createdAt := time.Now()
	twoHours, _ := time.ParseDuration("2h")
	expiresAt := createdAt.Add(twoHours)

	return &AuthToken{
		Id: id,
		DeviceId: deviceId,
		Token: token,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}, nil
}

func CheckAuth(db *sql.DB, token string) (*Device, error) {

	// check if token exists
	var deviceId int
	err := db.QueryRow(QUERY_TOKEN, token, time.Now()).Scan(&deviceId)
	if err != nil {
		return nil, err
	} 

	var uidString string
	var createdAt time.Time
	err = db.QueryRow(QUERY_DEVICE_BY_ID, deviceId).Scan(&uidString, &createdAt)
	if err != nil {
		return nil, err
	}

	uid, _ := uuid.FromString(uidString)

	return &Device{
		UUID: uid, 
		DeviceId: deviceId,
		CreatedAt: createdAt,
	}, nil
} 