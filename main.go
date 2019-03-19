package main

import (
	"github.com/moritzschramm/location-tracker-server/mqtt"
	"github.com/moritzschramm/location-tracker-server/api"
	"github.com/moritzschramm/location-tracker-server/config"
	"github.com/moritzschramm/location-tracker-server/database"
)

func main() {

	config := config.LoadConfig()

	db := database.SetupDatabase()
	defer db.Close()

	mqttClient := mqtt.SetupMQTTClient(db, config.MQTT)

	api.SetupAPI(db, mqttClient, config)
}
