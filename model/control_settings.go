package model

import (
	"database/sql"
	"time"
)

type ControlSettings struct {
	DB       *sql.DB   `json:"-"`
	Id       int       `json:"-"`
	DeviceId int       `json:"-"`
	UpdateAt time.Time `json:"updatedAt"`
}
