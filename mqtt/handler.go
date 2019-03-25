package mqtt

import (
	"strings"
	"strconv"
	"time"
	"database/sql"
	"log"
	"github.com/moritzschramm/location-tracker-server/model"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type SubHandler struct {
	DB *sql.DB
	Device *model.Device
	Client MQTT.Client
} 

func (handler *SubHandler) LocationCallback(client MQTT.Client, message MQTT.Message) {

	// message contains latitude and longitude in format <lat>,<long>,<time>
	m := strings.Split(string(message.Payload()), ",")

	lat, err := strconv.ParseFloat(m[0], 64)
	if err != nil {
		log.Println("Error parsing latitude: ", m[0], err.Error())
		return
	}
	long, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		log.Println("Error parsing longitude: ", m[1], err.Error())
		return
	}
	time, err := time.Parse(time.RFC3339, m[2])
	if err != nil {
		log.Println("Error parsing time: ", m[2], err.Error())
		return
	}

	_, err = model.MakeLocation(handler.DB, handler.Device.DeviceId, lat, long, time)
	if err != nil {
		log.Println("Error creating location: ", err.Error())
		return
	}
}

func (handler *SubHandler) BatteryCallback(client MQTT.Client, message MQTT.Message) {

}

func (handler *SubHandler) ControlSettingsCallback(client MQTT.Client, message MQTT.Message) {

}
