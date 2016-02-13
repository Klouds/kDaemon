package main

import (
	r "github.com/dancannon/gorethink"
	"github.com/superordinate/kDaemon/config"
	"github.com/superordinate/kDaemon/logging"
	"github.com/superordinate/kDaemon/routers"
	"github.com/superordinate/kDaemon/watcher"
	"net/http"
)

func main() {

	session, err := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "kdaemon",
	})

	if err != nil {
		//log.Panic(err.Error())
	}

	router := NewRouter(session)

	router.Handle("nodes subscribe", subscribeNodes)
	router.Handle("nodes unsubscribe", unsubscribeNodes)

	router.Handle("applications subscribe", subscribeApplications)
	router.Handle("applications unsubscribe", unsubscribeApplications)

	router.Handle("containers subscribe", subscribeContainers)
	router.Handle("containers unsubscribe", unsubscribeContainers)

	//Load the config file.
	err = config.LoadConfig()
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

	//Starts the cluster watcher
	go watcher.MainLoop()

	http.Handle("/", api.Mux)
	http.Handle("/ws", router)
	//Hosts the api server
	http.ListenAndServe(host+":"+apiport, nil)

}
