package model

import (
	"database/sql"
	"time"
)

const (
	LOCATION_INSERT_QUERY = "INSERT INTO locations (device_id, lat, long, time) VALUES (?, ?, ?, ?)"
)

type Location struct {
	DeviceId int       `json:"deviceId"`
	Lat      float64   `json:"lat"`
	Long     float64   `json:"long"`
	Time     time.Time `json:"time"`
}

func MakeLocation(db *sql.DB, deviceId int, lat, long float64, time time.Time) (*Location, error) {

	rows, err := db.Query(LOCATION_INSERT_QUERY, deviceId, lat, long, time)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return &Location{deviceId, lat, long, time}, nil
}

func GetLocations(db * sql.DB, from, to time.Time) ([]Location) {

	// TODO
	return nil
}
