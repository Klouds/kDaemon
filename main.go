package main

import (
	"net/http"
	"github.com/superordinate/kDaemon/routers"
	"github.com/superordinate/kDaemon/watcher"
	"github.com/superordinate/kDaemon/config"
	"github.com/superordinate/kDaemon/logging"

)

func main() {

	//Load the config file.
	err := config.LoadConfig()
	if err != nil {
		logging.Log("CONFIG FILE CANNOT BE LOADED")
		return
	}
	
	//Load some config file data
	host, err := config.Config.GetString("default", "bind_ip")
	if err != nil {
		logging.Log("Problem with config file! (bind_ip)")
		return
	}

	apiport, err := config.Config.GetString("default", "api_port")
	if err != nil {
		logging.Log("Problem with config file! (api_port)")
		return
	}

	uiport, err := config.Config.GetString("default", "ui_port")
	if err != nil {
		logging.Log("Problem with config file! (ui_port)")
		return
	}


	//Run the API
	var api routers.APIRouting
	api.Init()

	//Run the UI
	var ui routers.UIRouting
	ui.Init()
	
	//Starts the cluster watcher
	go watcher.MainLoop()

	//Hosts the api server
	go http.ListenAndServe(host + ":" + apiport, api.Mux)

	//Hosts the ui server
	http.ListenAndServe(host + ":" + uiport, ui.Mux)
	
}