package controllers

import (
	"github.com/julienschmidt/httprouter"
	"github.com/klouds/kDaemon/database"
	"github.com/klouds/kDaemon/models"
	"gopkg.in/unrolled/render.v1"
	"net/http"
)

type UIController struct {
	AppController
	*render.Render
}

func (c *UIController) Index(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	c.HTML(rw, http.StatusOK, "ui/index", nil)
}

//Nodes
//Landing page
func (c *UIController) NodeIndex(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	pageinfo := models.WebData{}
	//Check Login
	loggedIn := true

	pageinfo.LoggedIn = loggedIn

	if pageinfo.LoggedIn {

		nodes, err := database.GetNodes()

		if err != nil {
			nodes = nil
		}

		//placeholder user info
		pageinfo.CurrentUser = models.User{Username: "ozzadar"}
		pageinfo.PageData = nodes

		c.HTML(rw, http.StatusOK, "ui/nodes/index", pageinfo)
	} else {
		c.HTML(rw, http.StatusUnauthorized, "ui/notloggedin", pageinfo)
	}

}

//Create Node page
func (c *UIController) CreateNode(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	if r.Method == "GET" {
		pageinfo := models.WebData{}
		//Check Login
		loggedIn := true

		pageinfo.LoggedIn = loggedIn

		if pageinfo.LoggedIn {

			//placeholder user info
			pageinfo.CurrentUser = models.User{Username: "ozzadar"}
			pageinfo.PageData = nil

			c.HTML(rw, http.StatusOK, "ui/nodes/addnode", pageinfo)
		} else {
			c.HTML(rw, http.StatusUnauthorized, "ui/notloggedin", pageinfo)
		}
	}
}

//Containers

func (c *UIController) ContainerIndex(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	pageinfo := models.WebData{}
	//Check Login
	loggedIn := true

	pageinfo.LoggedIn = loggedIn

	if pageinfo.LoggedIn {

		conts, err := database.GetContainers()

		if err != nil {
			conts = nil
		}

		//placeholder user info
		pageinfo.CurrentUser = models.User{Username: "ozzadar"}
		pageinfo.PageData = conts

		c.HTML(rw, http.StatusOK, "ui/containers/index", pageinfo)
	} else {
		c.HTML(rw, http.StatusUnauthorized, "ui/notloggedin", pageinfo)
	}

}

//Applications

func (c *UIController) AppIndex(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	pageinfo := models.WebData{}
	//Check Login
	loggedIn := true

	pageinfo.LoggedIn = loggedIn

	if pageinfo.LoggedIn {

		apps, err := database.GetApplications()

		if err != nil {
			apps = nil
		}

		//placeholder user info
		pageinfo.CurrentUser = models.User{Username: "ozzadar"}
		pageinfo.PageData = apps

		c.HTML(rw, http.StatusOK, "ui/applications/index", pageinfo)
	} else {
		c.HTML(rw, http.StatusUnauthorized, "ui/notloggedin", pageinfo)
	}
}
