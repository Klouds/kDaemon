package main

import (
	"github.com/rs/cors"
	"github.com/superordinate/kDaemon/config"
	"github.com/superordinate/kDaemon/logging"
	"github.com/superordinate/kDaemon/routers"
	"github.com/superordinate/kDaemon/watcher"
	"net/http"
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

	//Run the API
	var api routers.APIRouting
	api.Init()

	var ws routers.WebSocketRouter
	ws.Init()

	//Starts the cluster watcher
	go watcher.MainLoop()

	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{"http://localhost:8081", "*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE"},
	})

	apihandler := c.Handler(api.Mux)
	wshandler := c.Handler(ws.Router)

	http.Handle("/", apihandler)
	http.Handle("/ws", wshandler)
	//Hosts the api server
	http.ListenAndServe(host+":"+apiport, nil)

}
