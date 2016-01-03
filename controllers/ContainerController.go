package controllers

import (
	"net/http"
	"gopkg.in/unrolled/render.v1"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/superordinate/kDaemon/watcher"
	"github.com/superordinate/kDaemon/models"
	//"github.com/superordinate/kDaemon/logging"
)

type ContainerController struct {
	AppController
	*render.Render
}

func (c *ContainerController) CreateContainer(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//creates a new application object populated with JSON from data
	newcontainer := models.Container{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newcontainer)

	if err != nil {
		panic(err)
		return
	}

	watcher.AddJob("LC", newcontainer)

	c.JSON(rw, http.StatusOK, newcontainer)
}

func (c *ContainerController) DeleteContainer(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Write([]byte("Removing Container: " + p.ByName("id")))
}

func (c *ContainerController) EditContainer(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Write([]byte("Editting Container: " + p.ByName("id")))
}

func (c *ContainerController) ContainerInformation(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Write([]byte("Showing info for Container: " + p.ByName("id")))
}
