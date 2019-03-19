package config

import (
	"log"

	"github.com/moritzschramm/location-tracker-server/mqtt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	MQTT      mqtt.Config `toml:"mqtt"`
	Host      string
	Port      string
	PublicDir string
}


func LoadConfig() Config {

	var config Config

	_, err := toml.DecodeFile("config.toml", &config)

	if err != nil {
		log.Panic(err)
	}

	return config
}
