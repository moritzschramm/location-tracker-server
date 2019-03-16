package main

import (
	"github.com/BurntSushi/toml"
)

var config Config

type Config struct {
	MQTT      MQTTConfig `toml:"mqtt"`
	Host      string
	Port      string
	PublicDir string
}

type MQTTConfig struct {
	Host     string
	Port     string
	CertFile string
	ClientId string
	Username string
	Password string
	PasswdFile string
}

func loadConfig() {

	_, err := toml.DecodeFile("config.toml", &config)

	if err != nil {
		panic(err)
	}
}
