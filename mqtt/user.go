package mqtt

import (
	"os/exec"
)

type User struct {
	Config Config
}

func (m User) AddUser(username, password string) error {

	cmd := exec.Command("mosquitto_passwd", "-b", m.Config.PasswdFile, username, password)

	return cmd.Run()
}

func (m User) DeleteUser(username string) error {

	cmd := exec.Command("mosquitto_passwd", "-D", m.Config.PasswdFile, username)

	return cmd.Run()
}