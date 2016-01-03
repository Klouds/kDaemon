package controllers

import (
	"net/http"
	"gopkg.in/unrolled/render.v1"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/superordinate/kDaemon/watcher"
	"github.com/superordinate/kDaemon/models"
	"github.com/superordinate/kDaemon/database"
	"strconv"
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


//This function must be passed as Jobs to the watcher, due to runtime container changes.
func (c *ContainerController) DeleteContainer(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Write([]byte("Removing Container: " + p.ByName("id")))
}

//This function must be passed as Jobs to the watcher, due to runtime container changes.
func (c *ContainerController) EditContainer(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Write([]byte("Editting Container: " + p.ByName("id")))
}

func (c *ContainerController) ContainerInformation(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//Gets the node id
	cont_id, err := strconv.Atoi(p.ByName("id"))

	if err != nil {
		c.JSON(rw, http.StatusBadRequest, "invalid id")
		return
	}

	//Attempts to retrieve the node from the database
	cont, err := database.GetContainer(int64(cont_id))

	if err != nil {
		c.JSON(rw, http.StatusNotFound, "Container doesn't exist")
		return
	}

	c.JSON(rw, http.StatusOK, cont)
}
