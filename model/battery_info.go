package model

import (
	"database/sql"
	"time"
)

const (
	INSERT_BATTERY_INFO = "INSERT INTO battery_infos (device_id, percentage, time) VALUES (?, ?, ?)"
	QUERY_BATTERY_INFO  = "SELECT * FROM battery_infos WHERE device_id = ? AND time >= ? AND time <= ?"
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
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &BatteryInfo{
		DB:         db,
		Id:         int(id),
		DeviceId:   deviceId,
		Percentage: percentage,
		Time:       time,
	}, nil
}

func GetBatteryInfo(db *sql.DB, deviceId int, from, to time.Time) ([]*BatteryInfo, error) {

	var infos []*BatteryInfo

	rows, err := db.Query(QUERY_BATTERY_INFO, deviceId, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		info := &BatteryInfo{}
		err = rows.Scan(&(info.Id), &(info.DeviceId), &(info.Percentage), &(info.Time))
		if err != nil {
			return nil, err
		}

		infos = append(infos, info)
	}

	return infos, nil
}
