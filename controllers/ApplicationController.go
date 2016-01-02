package controllers

import (
	"net/http"
	"gopkg.in/unrolled/render.v1"
	"github.com/julienschmidt/httprouter"
)

type ApplicationController struct {
	AppController
	*render.Render
}








func (c *ApplicationController) CreateApplication(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Write([]byte("Create Application"))
}

func (c *ApplicationController) DeleteApplication(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Write([]byte("Removing Application: " + p.ByName("id")))
}

func (c *ApplicationController) EditApplication(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Write([]byte("Editting Application: " + p.ByName("id")))
}

func (c *ApplicationController) ApplicationInformation(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Write([]byte("Showing info for Application: " + p.ByName("id")))
}
