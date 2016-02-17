package controllers

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/klouds/kDaemon/database"
	"github.com/klouds/kDaemon/models"
	"github.com/klouds/kDaemon/watcher"
	"gopkg.in/unrolled/render.v1"
	"net/http"
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

//This function must be passed as Jobs to the watcher, due to runtime container changes.
func (c *ContainerController) DeleteContainer(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	contid := p.ByName("id")

	oldcontainer, err := database.GetContainer(contid)

	if err != nil {
		c.JSON(rw, http.StatusNotFound, "Container doesn't exist")
		return
	}
	watcher.AddJob("RC", oldcontainer)

	c.JSON(rw, http.StatusOK, oldcontainer)
}

//This function must be passed as Jobs to the watcher, due to runtime container changes.
func (c *ContainerController) EditContainer(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Write([]byte("Editting Container: " + p.ByName("id")))
}

func (c *ContainerController) ContainerInformation(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//Gets the node id
	cont_id := p.ByName("id")

	//Attempts to retrieve the node from the database
	cont, err := database.GetContainer(cont_id)

	if err != nil {
		c.JSON(rw, http.StatusNotFound, "Container doesn't exist")
		return
	}

	c.JSON(rw, http.StatusOK, cont)
}

func (c *ContainerController) AllContainers(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	//Attempts to retrieve the node from the database
	conts, err := database.GetContainers()

	if err != nil {
		c.JSON(rw, http.StatusNotFound, "No containers")
		return
	}

	c.JSON(rw, http.StatusOK, conts)
}
