package main

import (
	"github.com/moritzschramm/location-tracker-server/api"
	"github.com/moritzschramm/location-tracker-server/config"
	"github.com/moritzschramm/location-tracker-server/database"
	"github.com/moritzschramm/location-tracker-server/mqtt"
)

func main() {

	config := config.Load()

	db := database.Setup(config)
	defer db.Close()

	mqttClient := mqtt.Setup(db, config.MQTT)

	api.Setup(db, mqttClient, config)
}
