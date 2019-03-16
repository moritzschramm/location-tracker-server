package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"io/ioutil"
	"log"
	"os"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func setupMQTTClient(db *sql.DB) *MQTT.Client {

	client := connectClient()

	setupMQTTSubscriptions(db, client)

	return client
}

func connectClient() {

	// set output of mqtt to stdout
	MQTT.DEBUG = log.New(os.Stdout, "", 0) // disable for production
	MQTT.ERROR = log.New(os.Stdout, "", 0)

	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(config.MQTT.CertFile)
	if err != nil {
		log.Println("Cannot read certificate")
		panic(err)
	}
	certpool.AppendCertsFromPEM(pemCerts)

	options := MQTT.NewClientOptions()
	options.AddBroker("ssl://" + config.MQTT.Host + ":" + config.MQTT.Port)
	options.SetTLSConfig(&tls.Config{RootCA: certpool})
	options.SetClientID(config.MQTT.ClientId)
	options.SetUsername(config.MQTT.Username)
	options.SetPassword(config.MQTT.Password)

	client := MQTT.NewClient(options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client
}

func setupMQTTSubscriptions(db *sql.DB, client *MQTT.Client) {

	var subCallback MQTT.MessageHandler = func(client MQTT.Client, message MQTT.message) {

		// do something here, e.g. store in database
	}

	if token := client.Subscribe("some/topic", 1, subCallback); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
	}
}
