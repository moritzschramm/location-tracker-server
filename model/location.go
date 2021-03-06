package model

import (
	"database/sql"
	"time"
)

const (
	INSERT_LOCATION = "INSERT INTO locations (device_id, lat, long, time) VALUES (?, ?, ?, ?)"
	QUERY_LOCATIONS = "SELECT * FROM locations WHERE device_id == ? AND time >= ? AND time <= ?"
)

type Location struct {
	Id       int       `json:"-"`
	DeviceId int       `json:"-"`
	Lat      float64   `json:"lat"`
	Long     float64   `json:"long"`
	Time     time.Time `json:"time"`
}

// create new location in database
func MakeLocation(db *sql.DB, deviceId int, lat, long float64, time time.Time) (*Location, error) {

	// insert location
	result, err := db.Exec(INSERT_LOCATION, deviceId, lat, long, time)
	if err != nil {
		return nil, err
	}

	// get id of location object
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Location{
		Id:       int(id),
		DeviceId: deviceId,
		Lat:      lat,
		Long:     long,
		Time:     time,
	}, nil
}

// return locations from database in specified time frame
func GetLocations(db *sql.DB, deviceId int, from, to time.Time) ([]*Location, error) {

	var locations []*Location

	// query all locations in time frame
	rows, err := db.Query(QUERY_LOCATIONS, deviceId, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// append each location to locations slice
	for rows.Next() {

		location := &Location{}
		err = rows.Scan(&(location.Id), &(location.DeviceId), &(location.Lat), &(location.Long), &(location.Time))
		if err != nil {
			return nil, err
		}

		locations = append(locations, location)
	}

	return locations, nil
}
