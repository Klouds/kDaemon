package config

import (
	"github.com/dlintw/goconf"
)

var Config *goconf.ConfigFile

func LoadConfig() error {

	err := error(nil)
	Config, err = goconf.ReadConfigFile("config/app.conf")

	return err

}
