package controllers

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/klouds/kDaemon/database"
	"github.com/klouds/kDaemon/models"
	"github.com/klouds/kDaemon/watcher2"
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

	//Removed the launcher call because
	//we are separating creation and launching.
	//we'll just add the container to the database

	database.CreateContainer(&newcontainer)
	c.JSON(rw, http.StatusOK, newcontainer)
}

//This function must be passed as Jobs to the watcher, due to runtime container changes.
func (c *ContainerController) LaunchContainer(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	contid := p.ByName("id")

	container, err := database.GetContainer(contid)

	if err != nil {
		c.JSON(rw, http.StatusNotFound, "Container doesn't exist")
		return
	}

	//this is where we would tell the server to launch the container
	//we would want to flag the container as 'launching'
	//and check whether it's launching before attempting another launch
	//
	//

	//Alright it's flagged as launching
	//now let's tell it to launch...
	//wait, reverse that order

	watcher2.TaskHandler.AddJob(watcher2.Launch,
		container.ApplicationID,
		container.Id,
		container.Name,
		"")

	//this is where we would tell the server to launch the container
	container.Status = "LAUNCHING"

	database.UpdateContainer(container)

	//I think that's all we need to launch the container
	//let's put that into our UI. and then get the backend cooperating

	//Actually... I think I'll add the launch button to the UI later.
	c.JSON(rw, http.StatusOK, container)

}

//This function must be passed as Jobs to the watcher, due to runtime container changes.
func (c *ContainerController) DeleteContainer(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	contid := p.ByName("id")

	err := database.DeleteContainer(contid)

	if err != nil {
		c.JSON(rw, http.StatusNotFound, "Container doesn't exist")
		return
	}

	c.JSON(rw, http.StatusOK, nil)
}

//This function must be passed as Jobs to the watcher, due to runtime container changes.
func (c *ContainerController) EditContainer(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//creates a new node object populated with JSON from data
	newcontainer := models.Container{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&newcontainer)
	if err != nil {
		panic(err)
		return
	}

	//Validates the Node passed in

	mergedcontainer, _ := database.GetContainer(p.ByName("id"))

	mergedcontainer = mergedcontainer.MergeChanges(&newcontainer)

	//Adds the node to the database
	success, _ := database.UpdateContainer(mergedcontainer)

	if success == false {
		c.JSON(rw, http.StatusNotFound, "Container doesn't exist")
		return
	}
	//return success message with new node information
	c.JSON(rw, http.StatusCreated, newcontainer)

}

func (c *ContainerController) StopContainer(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	contid := p.ByName("id")

	container, err := database.GetContainer(contid)

	if err != nil {
		c.JSON(rw, http.StatusNotFound, "Container doesn't exist")
		return
	}

	//this is where we would tell the server to launch the container
	//we would want to flag the container as 'launching'
	//and check whether it's launching before attempting another launch
	//
	//

	//Alright it's flagged as launching
	//now let's tell it to launch...
	//wait, reverse that order

	watcher2.TaskHandler.AddJob(watcher2.Stop,
		"",
		container.Id,
		container.Name,
		container.NodeID)

	//this is where we would tell the server to launch the container
	container.Status = "DOWN"

	database.UpdateContainer(container)
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
