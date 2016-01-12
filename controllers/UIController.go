package controllers

import (
	"net/http"
	"gopkg.in/unrolled/render.v1"
	"github.com/julienschmidt/httprouter"
)

type UIController struct {
	AppController
	*render.Render
}


func (c *UIController) Index(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Write([]byte("Index"))
}
