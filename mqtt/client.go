package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"io/ioutil"
	"log"
	"os"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func SetupMQTTClient(db *sql.DB, config Config) MQTT.Client {

	client := connectClient(config)

	setupMQTTSubscriptions(db, client)

	return client
}

func connectClient(config Config) MQTT.Client {

	// set output of mqtt to stdout
	MQTT.DEBUG = log.New(os.Stdout, "", 0) // disable for production
	MQTT.ERROR = log.New(os.Stdout, "", 0)

	options := MQTT.NewClientOptions()
	options.AddBroker("ssl://" + config.Host + ":" + config.Port)
	options.SetTLSConfig(TLSConfig(config))
	options.SetClientID(config.ClientId)
	options.SetUsername(config.Username)
	options.SetPassword(config.Password)

	client := MQTT.NewClient(options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Panic(token.Error())
	}

	return client
}

func TLSConfig(config Config) *tls.Config {

	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(config.CertFile)
	if err != nil {
		log.Panic("Cannot read certificate", err.Error())
	}
	certpool.AppendCertsFromPEM(pemCerts)

	return &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		InsecureSkipVerify: true,
		ClientCAs:          nil,
		Certificates:       nil,
	}
}

func setupMQTTSubscriptions(db *sql.DB, client MQTT.Client) {

	var subCallback MQTT.MessageHandler = func(client MQTT.Client, message MQTT.Message) {

		// do something here, e.g. store in database
	}

	if token := client.Subscribe("some/topic", 1, subCallback); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
	}
}
