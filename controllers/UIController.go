package controllers

import (
	"github.com/julienschmidt/httprouter"
	"gopkg.in/unrolled/render.v1"
	"net/http"
)

type UIController struct {
	AppController
	*render.Render
}

func (c *UIController) Index(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Write([]byte("Index"))
}

//Nodes

//Containers

//Applications
