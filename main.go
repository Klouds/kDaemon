package main

import (
	"net/http"
	//"fmt"
	"github.com/superordinate/kDaemon/routers"
)

func main() {

	var newmux routers.Routing
	newmux.Init()
	
	http.ListenAndServe("0.0.0.0:1337", newmux.Mux)
}