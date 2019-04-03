package config

import (
	"log"

	"github.com/moritzschramm/location-tracker-server/mqtt"

	"github.com/BurntSushi/toml"
)

const (
	CONFIG_FILE = "config/config.toml"
)

type Config struct {
	MQTT          mqtt.Config `toml:"mqtt"`
	Host          string
	Port          string
	PublicDir     string
	AdminUUID     string
	AdminPassword string
	CertFile      string
	KeyFile       string
}

func LoadConfig() Config {

	var config Config

	_, err := toml.DecodeFile(CONFIG_FILE, &config)

	if err != nil {
		log.Panic(err)
	}

	return config
}
