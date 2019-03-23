package model

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	QUERY_TOKEN          = "SELECT id, device_id, created_at, expires_at FROM tokens WHERE token = '?' AND expires_at <= ?"
	QUERY_DEVICE_BY_UUID = "SELECT device_id, password FROM devices WHERE uuid == ?"
	INSERT_TOKEN         = "INSERT INTO tokens (device_id, token, created_at, expires_at) VALUES (?, ?, ?, ?)"
	DELETE_TOKEN         = "DELETE FROM tokens WHERE id == ?"
)

type AuthToken struct {
	DB        *sql.DB   `json:"-"`
	Id        int       `json:"-"`
	DeviceId  int       `json:"-"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"-"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func (token *AuthToken) Logout() error {

	_, err := token.DB.Exec(DELETE_TOKEN, token.Id)
	if err != nil {
		return err
	}

	token = nil

	return nil
}

func (token *AuthToken) Refresh() (*AuthToken, error) {

	db := token.DB
	deviceId := token.DeviceId

	// delete old token
	token.Logout()

	return createNewToken(db, deviceId)
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

	return createNewToken(db, deviceId)
}

func GetAuthToken(db *sql.DB, token string) (*AuthToken, error) {

	var id int
	var deviceId int
	var createdAt time.Time
	var expiresAt time.Time
	err := db.QueryRow(QUERY_TOKEN, token, time.Now()).Scan(&id, &deviceId, &createdAt, &expiresAt)
	if err != nil {
		return nil, err
	}

	return &AuthToken{
		Id:        id,
		DeviceId:  deviceId,
		Token:     token,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}, nil
}

func createNewToken(db *sql.DB, deviceId int) (*AuthToken, error) {

	// create new random token (256 bit long)
	// TODO
	tokenString := "hello"
	createdAt := time.Now()
	expiresAt := createdAt.Add(2 * time.Hour)

	token := &AuthToken{
		DB:        db,
		Id:        0,
		DeviceId:  deviceId,
		Token:     tokenString,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}

	// insert new token in database
	// TODO
	var id int

	token.Id = id

	return token, nil
}
