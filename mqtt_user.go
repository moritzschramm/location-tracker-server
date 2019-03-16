package main

import (
	"os/exec"
)

type MQTTUser struct {
	Config MQTTConfig
}

func (m MQTTUser) AddUser(username, password string) error {

	cmd := exec.Command("mosquitto_passwd", "-b", m.Config.PasswdFile, username, password)

	return cmd.Run()
}

func (m MQTTUser) DeleteUser(username string) error {

	cmd := exec.Command("mosquitto_passwd", "-D", m.Config.PasswdFile, username)

	return cmd.Run()
}