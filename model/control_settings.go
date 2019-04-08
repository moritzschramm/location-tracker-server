package model

import (
	"database/sql"
	"time"
)

const (
	QUERY_CONTROL_SETTINGS  = "SELECT * FROM control_settings WHERE device_id = ?"
	INSERT_CONTROL_SETTINGS = "INSERT INTO control_settings (device_id, updated_at, operation_mode, alarm, alarm_enabled, update_frequency, rf_enabled VALUES (?, ?, ?, ?, ?, ?, ?)"
	UPDATE_CONTROL_SETTINGS = "UPDATE control_settings SET updated_at=?, oparation_mode=?, alarm=?, update_frequency=?, rf_enabled=? WHERE device_id = ?"

	OP_15 = 0
)

type ControlSettings struct {
	DB        *sql.DB   `json:"-"`
	DeviceId  int       `json:"-"`
	UpdatedAt time.Time `json:"updatedAt"`

	OperationMode   int  `json:"operationMode"`
	Alarm           bool `json:"alarm"`
	AlarmEnabled    bool `json:"alarmEnabled"`
	UpdateFrequency int  `json:"updateFrequency"`
	RFEnabled       bool `json:"RFEnabled"`
}

// update control settings in database
// insert new control settings if not exists
func (settings *ControlSettings) Update() error {

	err := settings.DB.QueryRow(QUERY_CONTROL_SETTINGS, settings.DeviceId).Scan()

	if err != nil { // no control setting for this device in database, insert new settings

		_, err := settings.DB.Exec(INSERT_CONTROL_SETTINGS,
			settings.DeviceId,
			settings.UpdatedAt,
			settings.OperationMode,
			settings.Alarm,
			settings.AlarmEnabled,
			settings.UpdateFrequency,
			settings.RFEnabled)

		if err != nil {
			return err
		}

	} else {

		_, err := settings.DB.Exec(UPDATE_CONTROL_SETTINGS,
			settings.UpdatedAt,
			settings.OperationMode,
			settings.Alarm,
			settings.AlarmEnabled,
			settings.UpdateFrequency,
			settings.RFEnabled,
			settings.DeviceId)

		if err != nil {
			return err
		}
	}

	return nil
}
