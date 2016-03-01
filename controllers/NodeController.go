package controllers

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/klouds/kDaemon/database"
	"github.com/klouds/kDaemon/logging"
	"github.com/klouds/kDaemon/models"
	"github.com/klouds/kDaemon/watcher"
	"gopkg.in/unrolled/render.v1"
	"net/http"
)

type NodeController struct {
	AppController
	*render.Render
}

func (c *NodeController) CreateNode(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//creates a new node object populated with JSON from data
	newnode := models.Node{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newnode)
	if err != nil {
		panic(err)
		return
	}

	//Validates the Node passed in

	if newnode.Validate() {
		//Adds the node to the database
		nodeid, err := database.CreateNode(&newnode)

		if err != nil {
			c.JSON(rw, http.StatusConflict, "Node conflicts with existing node. Make sure your node is unique.")
			return
		}

		//We're going to add a "add node" job to the new watcher

		watcher.TaskHandler.AddJob(watcher.AddNode, "", "", "", nodeid)
		//return success message with new node information
		c.JSON(rw, http.StatusCreated, newnode)
	} else {

		body, _ := newnode.GetJSON()
		c.JSON(rw, http.StatusBadRequest, "Invalid format"+body)
	}

}

func (c *NodeController) DeleteNode(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//Gets the node id
	nodeid := p.ByName("id")

	//Attempts to remove the node
	success, _ := database.DeleteNode(nodeid)

	if !success {
		c.JSON(rw, http.StatusNotFound, "Node doesn't exist")
		return
	}

	c.JSON(rw, http.StatusOK, "Node deleted successfully")

}

func (c *NodeController) EditNode(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	//creates a new node object populated with JSON from data
	newnode := models.Node{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&newnode)
	if err != nil {
		panic(err)
		return
	}

	//Validates the Node passed in

	mergedNode, _ := database.GetNode(p.ByName("id"))

	mergedNode = mergedNode.MergeChanges(&newnode)

	logging.Log(newnode.ContainerCount)

	if mergedNode.Validate() {
		//Adds the node to the database
		success, _ := database.UpdateNode(mergedNode)

		if success == false {
			c.JSON(rw, http.StatusNotFound, "Node doesn't exist")
			return
		}
		//return success message with new node information
		c.JSON(rw, http.StatusCreated, newnode)
	} else {
		c.JSON(rw, http.StatusBadRequest, "Invalid format")
	}
}

func (c *NodeController) NodeInformation(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	//Gets the node id
	nodeid := p.ByName("id")

	//Attempts to retrieve the node from the database
	node, err := database.GetNode(nodeid)

	if err != nil {
		c.JSON(rw, http.StatusNotFound, "Node doesn't exist")
		return
	}

	c.JSON(rw, http.StatusOK, node)
}

func (c *NodeController) AllNodes(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	//Attempts to retrieve the node from the database
	nodes, err := database.GetNodes()

	if err != nil {
		c.JSON(rw, http.StatusNotFound, "No nodes")
		return
	}

	c.JSON(rw, http.StatusOK, nodes)
}
