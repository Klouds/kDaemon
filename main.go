package main

import (
	"net/http"
	//"fmt"
	"github.com/superordinate/kDaemon/routers"
	"github.com/superordinate/kDaemon/watcher"
)

func main() {

	var newmux routers.Routing
	newmux.Init()
	
	go watcher.MainLoop()
	http.ListenAndServe("0.0.0.0:1337", newmux.Mux)
	
	
}