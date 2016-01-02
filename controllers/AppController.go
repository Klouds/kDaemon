package controllers

import (
	"net/http"
)
type Action func (rw http.ResponseWriter, r * http.Request) error

//This is our base controller
type AppController struct {}

//The Action functions helps with error handling in a controller
func (c *AppController) Action(a Action) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if err := a(rw, r); err != nil {
				http.Error(rw, err.Error(), 500)
			}

		})

}