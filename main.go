package main

import (
	"net/http"
	"github.com/superordinate/kDaemon/routers"
	"github.com/superordinate/kDaemon/watcher"
)

func main() {

	var newmux routers.Routing
	newmux.Init()
	
	//Starts the cluster watcher
	go watcher.MainLoop()

	//Hosts the web server
	http.ListenAndServe("0.0.0.0:1337", newmux.Mux)
	
	
}