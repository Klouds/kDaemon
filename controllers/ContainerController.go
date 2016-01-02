package controllers

import (
	"net/http"
	"gopkg.in/unrolled/render.v1"
	"github.com/julienschmidt/httprouter"
)

type ContainerController struct {
	AppController
	*render.Render
}








func (c *ContainerController) CreateContainer(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Write([]byte("Create Container"))
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
