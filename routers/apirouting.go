package routers

import (
	"github.com/julienschmidt/httprouter"
	"github.com/klouds/kDaemon/config"
	"github.com/klouds/kDaemon/controllers"
	"github.com/klouds/kDaemon/database"
	"github.com/klouds/kDaemon/logging"
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

	APIVERSION, err := config.Config.GetString("default", "api_version")
	if err != nil {
		logging.Log("Problem with config file! (api_version)")
	}
	/*Basics of controller*/
	ac := &controllers.ApplicationController{Render: r.Render}
	cc := &controllers.ContainerController{Render: r.Render}
	nc := &controllers.NodeController{Render: r.Render}
	uc := &controllers.UserController{Render: r.Render}

	//Application Management

	r.Mux.POST("/"+APIVERSION+"/applications/create", ac.CreateApplication)
	r.Mux.DELETE("/"+APIVERSION+"/applications/delete/:id", ac.DeleteApplication)
	r.Mux.PATCH("/"+APIVERSION+"/applications/update/:id", ac.EditApplication)
	r.Mux.GET("/"+APIVERSION+"/applications/:id", ac.ApplicationInformation)
	r.Mux.GET("/"+APIVERSION+"/applications/:id/", ac.ApplicationInformation)
	r.Mux.GET("/"+APIVERSION+"/applications", ac.AllApplications)

	//Container Management
	r.Mux.POST("/"+APIVERSION+"/containers/create", cc.CreateContainer)
	r.Mux.POST("/"+APIVERSION+"/containers/launch/:id", cc.LaunchContainer)
	r.Mux.POST("/"+APIVERSION+"/containers/stop/:id", cc.StopContainer)
	r.Mux.DELETE("/"+APIVERSION+"/containers/delete/:id", cc.DeleteContainer)
	r.Mux.PATCH("/"+APIVERSION+"/containers/update/:id", cc.EditContainer)
	r.Mux.GET("/"+APIVERSION+"/containers/:id", cc.ContainerInformation)
	r.Mux.GET("/"+APIVERSION+"/containers/:id/", cc.ContainerInformation)
	r.Mux.GET("/"+APIVERSION+"/containers", cc.AllContainers)

	//Node Management
	r.Mux.POST("/"+APIVERSION+"/nodes/create", nc.CreateNode)
	r.Mux.DELETE("/"+APIVERSION+"/nodes/delete/:id", nc.DeleteNode)
	r.Mux.PATCH("/"+APIVERSION+"/nodes/update/:id", nc.EditNode)
	r.Mux.GET("/"+APIVERSION+"/nodes/:id", nc.NodeInformation)
	r.Mux.GET("/"+APIVERSION+"/nodes/:id/", nc.NodeInformation)
	r.Mux.GET("/"+APIVERSION+"/nodes", nc.AllNodes)

	//User Management
	r.Mux.POST("/"+APIVERSION+"/users/create", uc.CreateUser)
	r.Mux.DELETE("/"+APIVERSION+"/users/delete/:id", uc.DeleteUser)
	r.Mux.PATCH("/"+APIVERSION+"/users/update/:id", uc.EditUser)
	r.Mux.GET("/"+APIVERSION+"/users/:id", uc.UserInformation)
	r.Mux.GET("/"+APIVERSION+"/users/:id/", uc.UserInformation)

	logging.Log("Web server active")
}
