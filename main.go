package main

import (
	"github.com/klouds/kDaemon/config"
	"github.com/klouds/kDaemon/logging"
	"github.com/klouds/kDaemon/routers"
	//"github.com/klouds/kDaemon/watcher"
	"github.com/klouds/kDaemon/watcher2"
	"github.com/rs/cors"
	"net/http"
	// "time"
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
	//go watcher.MainLoop()

	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{"http://localhost:8081", "*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE"},
	})

	watcher_new := watcher2.Watcher{}
	watcher_new.Init()

	stop := make(chan bool)

	go watcher_new.Run(stop)

	apihandler := c.Handler(api.Mux)
	wshandler := c.Handler(ws.Router)

	http.Handle("/", apihandler)
	http.Handle("/ws", wshandler)
	//Hosts the api server

	// time.Sleep(10 * time.Second)
	// stop <- true

	http.ListenAndServe(host+":"+apiport, nil)

}
