package controllers

import (
	//"fmt"
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

	//Adds the node to the database


	//return success message with new node information
	c.JSON(rw, http.StatusOK, newnode)
}

func (c *NodeController) DeleteNode(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Write([]byte("Removing Node: " + p.ByName("id")))
}

func (c *NodeController) EditNode(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Write([]byte("Editting Node: " + p.ByName("id")))
}

func (c *NodeController) NodeInformation(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	/*
	node := &models.Node{
		Id: 0,
		UserID: 0,
		Hostname: "testNode",
		DIPAddr: "127.0.0.1",
		DPort: "2575",
		PIPAddr: "127.0.0.1",
		PPort: "9090",
		IsEnabled: true,
	}
*/
	rw.Write([]byte("Showing info for Node: " + p.ByName("id")))
}
