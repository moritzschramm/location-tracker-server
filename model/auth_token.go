package model

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

const (
	SECURE_COOKIE = false // TODO enable for production
	TOKEN_SIZE    = 256   // token size in bits

	QUERY_TOKEN          = "SELECT id, device_id, created_at, expires_at FROM tokens WHERE token = ? AND expires_at >= ?"
	QUERY_DEVICE_BY_UUID = "SELECT device_id, password FROM devices WHERE uuid = ?"
	INSERT_TOKEN         = "INSERT INTO tokens (device_id, token, created_at, expires_at) VALUES (?, ?, ?, ?)"
	DELETE_TOKEN         = "DELETE FROM tokens WHERE id = ?"
)

type AuthToken struct {
	DB        *sql.DB   `json:"-"`
	Id        int       `json:"-"`
	DeviceId  int       `json:"-"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"-"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// return max age in seconds until auth token expires
func (token *AuthToken) MaxAge() int {

	return int(time.Until(token.ExpiresAt).Seconds())
}

// create cookie from auth token
func (token *AuthToken) ToCookie() *http.Cookie {

	return &http.Cookie{
		Name:     "token",
		Value:    token.Token,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Expires:  token.ExpiresAt,
		Secure:   SECURE_COOKIE,
		HttpOnly: true,
	}
}

// create (empty) cookie to unset token
func (token *AuthToken) UnsetCookie() *http.Cookie {

	return &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
		Secure:   SECURE_COOKIE,
		HttpOnly: true,
	}
}

// make auth token invalid
func (token *AuthToken) Logout() error {

	_, err := token.DB.Exec(DELETE_TOKEN, token.Id)
	if err != nil {
		return err
	}

	return nil
}

// refresh auth token
// old token will be invalid, new token will be returned
func (token *AuthToken) Refresh() (*AuthToken, error) {

	db := token.DB
	deviceId := token.DeviceId

	// delete old token
	token.Logout()

	return createNewToken(db, deviceId)
}

// authenticate device (check UUID and password)
// return auth token if authentication succesful
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

// get auth token from database by token
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
		DB:        db,
		Id:        id,
		DeviceId:  deviceId,
		Token:     token,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}, nil
}

// create random token, create and insert auth token into database
func createNewToken(db *sql.DB, deviceId int) (*AuthToken, error) {

	// create new random token (256 bit long)
	tokenString, err := generateRandomToken(TOKEN_SIZE)
	if err != nil {
		return nil, err
	}

	createdAt := time.Now()
	expiresAt := createdAt.Add(2 * time.Hour)

	// insert new token in database
	result, err := db.Exec(INSERT_TOKEN, deviceId, tokenString, createdAt, expiresAt)
	if err != nil {
		return nil, err
	}

	// get id of auth token
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	token := &AuthToken{
		DB:        db,
		Id:        int(id),
		DeviceId:  deviceId,
		Token:     tokenString,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}

	return token, nil
}

// helper function to generate random token
func generateRandomToken(bytes int) (string, error) {

	b := make([]byte, bytes)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	s := base64.StdEncoding.EncodeToString(b)

	return s, nil
}
