package main

import (
	"net/http"
	"github.com/superordinate/kDaemon/routers"
	"github.com/superordinate/kDaemon/watcher"
	"github.com/superordinate/kDaemon/config"
	"github.com/superordinate/kDaemon/logging"

)

func main() {

	err := config.LoadConfig()

	if err != nil {
		logging.Log("CONFIG FILE CANNOT BE LOADED")
		return
	}

	logging.Log(config.Config)
	
	host, err := config.Config.GetString("default", "bind_ip")

	if err != nil {
		logging.Log("Problem with config file! (bind_ip)")
	}

	port, err := config.Config.GetString("default", "bind_port")
	if err != nil {
		logging.Log("Problem with config file! (bind_port)")
	}

	var newmux routers.Routing
	newmux.Init()
	
	//Starts the cluster watcher
	go watcher.MainLoop()

	//Hosts the web server
	http.ListenAndServe(host + ":" + port, newmux.Mux)
	
	
}