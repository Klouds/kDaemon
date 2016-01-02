package controllers

import (
	"net/http"
	"gopkg.in/unrolled/render.v1"
	"github.com/julienschmidt/httprouter"
)

type UserController struct {
	AppController
	*render.Render
}


func (c *UserController) CreateUser(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Write([]byte("Create User"))
}

func (c *UserController) DeleteUser(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Write([]byte("Removing User: " + p.ByName("id")))
}

func (c *UserController) EditUser(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Write([]byte("Editting User: " + p.ByName("id")))
}

func (c *UserController) UserInformation(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Write([]byte("Showing info for User: " + p.ByName("id")))
}
