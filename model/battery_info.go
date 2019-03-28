package model

import (
	"database/sql"
	"time"
)

const (
	INSERT_BATTERY_INFO = "INSERT INTO battery_infos (device_id, percentage, time) VALUES (?, ?, ?)"
)

type BatteryInfo struct {
	DB         *sql.DB   `json:"-"`
	Id         int       `json:"-"`
	DeviceId   int       `json:"-"`
	Percentage int       `json:"percentage"`
	Time       time.Time `json:"time"`
}

func MakeBatteryInfo(db *sql.DB, deviceId, percentage int, time time.Time) (*BatteryInfo, error) {

	result, err := db.Exec(INSERT_BATTERY_INFO, deviceId, percentage, time)
	if err != nil {
		return nil, error
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, error
	}

	return &BatteryInfo{
		DB:         db,
		Id:         int(id),
		DeviceId:   deviceId,
		Percentage: percentage,
		Time:       time,
	}, nil
}
