package mqtt

import (
	"os/exec"
)

type Config struct {
	Host       string
	Port       string
	CertFile   string
	ClientId   string
	Username   string
	Password   string
	PasswdFile string
}

// add user to mosquitto passwd file
func (c Config) AddUser(username, password string) error {

	cmd := exec.Command("mosquitto_passwd", "-b", c.PasswdFile, username, password)

	return cmd.Run()
}

// remove user from mosquitto passwd file
func (c Config) DeleteUser(username string) error {

	cmd := exec.Command("mosquitto_passwd", "-D", c.PasswdFile, username)

	return cmd.Run()
}
