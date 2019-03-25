package model

import (
	"database/sql"
	"time"
)

type BatteryInfo struct {
	DB         *sql.DB   `json:"-"`
	Id         int       `json:"-"`
	DeviceId   int       `json:"-"`
	Percentage int       `json:"percentage"`
	Time       time.Time `json:"time"`
}

func MakeBatteryInfo(db *sql.DB, deviceId, percentage int, time time.Time) (*BatteryInfo, error) {

	return nil, nil
}
