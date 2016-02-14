package controllers

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/superordinate/kDaemon/database"
	"github.com/superordinate/kDaemon/models"
	"gopkg.in/unrolled/render.v1"
	"net/http"
)

type ApplicationController struct {
	AppController
	*render.Render
}

func (c *ApplicationController) CreateApplication(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	//creates a new application object populated with JSON from data
	newapp := models.Application{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newapp)

	if err != nil {
		panic(err)
		return
	}

	//Validates the Node passed in

	if newapp.Validate() {
		//Adds the node to the database
		success, _ := database.CreateApplication(&newapp)

		if success == false {
			c.JSON(rw, http.StatusConflict, "Application conflict. Make sure your application is unique.")
			return
		}
		//return success message with new node information
		c.JSON(rw, http.StatusCreated, newapp)
	} else {
		c.JSON(rw, http.StatusBadRequest, newapp)
	}
}

func (c *ApplicationController) DeleteApplication(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//Gets the app id
	appid := p.ByName("id")

	//Attempts to remove the node
	success, _ := database.DeleteApplication(appid)

	if !success {
		c.JSON(rw, http.StatusNotFound, "Application doesn't exist")
		return
	}

	c.JSON(rw, http.StatusOK, "Application deleted successfully")
}

func (c *ApplicationController) EditApplication(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//creates a new application object populated with JSON from data
	app := models.Application{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&app)
	if err != nil {
		panic(err)
		return
	}

	app.Id = p.ByName("id")
	//Validates the Node passed in

	if app.Validate() {
		//Adds the node to the database
		success, _ := database.UpdateApplication(&app)

		if success == false {
			c.JSON(rw, http.StatusNotFound, "Application doesn't exist")
			return
		}
		//return success message with new node information
		c.JSON(rw, http.StatusCreated, app)
	} else {
		c.JSON(rw, http.StatusBadRequest, "Invalid format")
	}
}

func (c *ApplicationController) ApplicationInformation(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//Gets the app id
	appid := p.ByName("id")

	//Attempts to retrieve the application from the database
	app, err := database.GetApplication(appid)

	if err != nil {
		c.JSON(rw, http.StatusNotFound, "Node doesn't exist")
		return
	}

	c.JSON(rw, http.StatusOK, app)
}

func (c *ApplicationController) AllApplications(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	//Attempts to retrieve the node from the database
	apps, err := database.GetApplications()

	if err != nil {
		c.JSON(rw, http.StatusNotFound, "No apps")
		return
	}

	c.JSON(rw, http.StatusOK, apps)
}
