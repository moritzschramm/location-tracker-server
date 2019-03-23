package mqtt

import (
	"github.com/moritzschramm/location-tracker-server/model"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type SubCallback func(MQTT.Client, MQTT.Message, *model.Device)

func LocationCallback(client MQTT.Client, message MQTT.Message, device *model.Device) {

}

func BatteryCallback(client MQTT.Client, message MQTT.Message, device *model.Device) {

}

func ControlSettingsCallback(client MQTT.Client, message MQTT.Message, device *model.Device) {

}
