package main

func main() {

	loadConfig()

	db := setupDatabase()
	defer db.Close()

	mqttClient := setupMQTTClient(db)

	setupAPI(db, mqttClient)
}
