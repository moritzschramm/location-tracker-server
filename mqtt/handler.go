package mqtt

import (
	"database/sql"
	"github.com/moritzschramm/location-tracker-server/model"
	"log"
	"strconv"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type SubHandler struct {
	DB     *sql.DB
	Device *model.Device
	Client MQTT.Client
}

func (handler *SubHandler) SubscribeTo(topic string, callback MQTT.MessageHandler) {

	token := handler.Client.Subscribe(handler.Device.UUID.String()+topic, 1, callback)
	if token.Wait() && token.Error() != nil {
		log.Println(token.Error())
	}
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
		log.Println("Error parsing location time: ", m[2], err.Error())
		return
	}

	_, err = model.MakeLocation(handler.DB, handler.Device.DeviceId, lat, long, time)
	if err != nil {
		log.Println("Error creating location: ", err.Error())
		return
	}
}

func (handler *SubHandler) BatteryInfoCallback(client MQTT.Client, message MQTT.Message) {

	// message contains battery info in format <percentage>,<time>
	m := strings.Split(string(message.Payload()), ",")

	percentage, err := strconv.ParseInt(m[0], 10, 32)
	if err != nil {
		log.Println("Error parsing percentage: ", m[0], err.Error())
		return
	}

	time, err := time.Parse(time.RFC3339, m[1])
	if err != nil {
		log.Println("Error parsing battery info time: ", m[1], err.Error())
		return
	}

	_, err = model.MakeBatteryInfo(handler.DB, handler.Device.DeviceId, int(percentage), time)
	if err != nil {
		log.Println("Error creating battery info: ", err.Error())
		return
	}
}

func (handler *SubHandler) ControlSettingsCallback(client MQTT.Client, message MQTT.Message) {

	// message contains battery info in format <opMode>,<alarmEnabled>,<updateFrequency>,<RFEnabled>
	//m := strings.Split(string(message.Payload()), ",")

	// TODO
}

func (handler *SubHandler) AlarmCallback(client MQTT.Client, message MQTT.Message) {

}
