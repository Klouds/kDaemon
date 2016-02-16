package routers

import (
	"github.com/julienschmidt/httprouter"
	"github.com/superordinate/kDaemon/controllers"
	"github.com/superordinate/kDaemon/database"
	"github.com/superordinate/kDaemon/logging"
	"gopkg.in/unrolled/render.v1"
	"html/template"
)

type APIRouting struct {
	Render *render.Render
	Mux    *httprouter.Router
}

func (r *APIRouting) Init() {

	database.Init()
	r.Render = render.New(render.Options{Directory: "views",
		IndentJSON: true,
		Funcs: []template.FuncMap{
			{

				"str2html": func(raw string) template.HTML {
					return template.HTML(raw)
				},
				"add": func(x, y int) int {
					return x + y
				},
				"mod": func(x, y int) int {
					return x % y
				},
			},
		},
	})
	r.Mux = httprouter.New()

	APIVERSION := "0.0"
	/*Basics of controller*/
	ac := &controllers.ApplicationController{Render: r.Render}
	cc := &controllers.ContainerController{Render: r.Render}
	nc := &controllers.NodeController{Render: r.Render}
	uc := &controllers.UserController{Render: r.Render}

	//Application Management

	r.Mux.POST("/"+APIVERSION+"/applications/create", ac.CreateApplication)
	r.Mux.DELETE("/"+APIVERSION+"/applications/:id/delete", ac.DeleteApplication)
	r.Mux.PATCH("/"+APIVERSION+"/applications/:id/update", ac.EditApplication)
	r.Mux.GET("/"+APIVERSION+"/applications/:id", ac.ApplicationInformation)
	r.Mux.GET("/"+APIVERSION+"/applications/:id/", ac.ApplicationInformation)
	r.Mux.GET("/"+APIVERSION+"/applications", ac.AllApplications)

	//Container Management
	r.Mux.POST("/"+APIVERSION+"/containers/create", cc.CreateContainer)
	r.Mux.DELETE("/"+APIVERSION+"/containers/:id/delete", cc.DeleteContainer)
	r.Mux.PATCH("/"+APIVERSION+"/containers/:id/update", cc.EditContainer)
	r.Mux.GET("/"+APIVERSION+"/containers/:id", cc.ContainerInformation)
	r.Mux.GET("/"+APIVERSION+"/containers/:id/", cc.ContainerInformation)
	r.Mux.GET("/"+APIVERSION+"/containers", cc.AllContainers)

	//Node Management
	r.Mux.POST("/"+APIVERSION+"/nodes/create", nc.CreateNode)
	r.Mux.DELETE("/"+APIVERSION+"/nodes/:id/delete", nc.DeleteNode)
	r.Mux.PATCH("/"+APIVERSION+"/nodes/:id/update", nc.EditNode)
	r.Mux.GET("/"+APIVERSION+"/nodes/:id", nc.NodeInformation)
	r.Mux.GET("/"+APIVERSION+"/nodes/:id/", nc.NodeInformation)
	r.Mux.GET("/"+APIVERSION+"/nodes", nc.AllNodes)

	//User Management
	r.Mux.POST("/"+APIVERSION+"/users/create", uc.CreateUser)
	r.Mux.DELETE("/"+APIVERSION+"/users/:id/delete", uc.DeleteUser)
	r.Mux.PATCH("/"+APIVERSION+"/users/:id/update", uc.EditUser)
	r.Mux.GET("/"+APIVERSION+"/users/:id", uc.UserInformation)
	r.Mux.GET("/"+APIVERSION+"/users/:id/", uc.UserInformation)

	logging.Log("Web server active")
}
