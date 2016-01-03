package controllers

import (
	"strconv"
	"net/http"
	"gopkg.in/unrolled/render.v1"
	"encoding/json"
	"github.com/superordinate/kDaemon/models"
	"github.com/julienschmidt/httprouter"
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
		success, _ := CreateNode(&newnode)

		if success == false {
			c.JSON(rw, http.StatusConflict, "Node conflicts with existing node. Make sure your node is unique.")
			return
		}
		//return success message with new node information
		c.JSON(rw, http.StatusCreated, newnode)
	} else {
		c.JSON(rw, http.StatusBadRequest, "Invalid format")
	}
	
}

func (c *NodeController) DeleteNode(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//Gets the node id
	nodeid, err := strconv.Atoi(p.ByName("id"))

	if err != nil {
		c.JSON(rw, http.StatusBadRequest, "invalid id")
		return
	}

	//Attempts to remove the node
	success, _ := DeleteNode(int64(nodeid))

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

	if newnode.Validate() {
		//Adds the node to the database
		success, _ := UpdateNode(&newnode)

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
	nodeid, err := strconv.Atoi(p.ByName("id"))

	if err != nil {
		c.JSON(rw, http.StatusBadRequest, "invalid id")
		return
	}

	//Attempts to retrieve the node from the database
	node, err := GetNode(int64(nodeid))

	if err != nil {
		c.JSON(rw, http.StatusNotFound, "Node doesn't exist")
		return
	}

	c.JSON(rw, http.StatusOK, node)
}
