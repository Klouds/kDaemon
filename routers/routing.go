package routers

import (
	"net/http"
	"html/template"
	"gopkg.in/unrolled/render.v1"
	"github.com/julienschmidt/httprouter"
	"github.com/superordinate/kDaemon/controllers"
	"fmt"
)

type Routing struct {

	Render *render.Render
	Mux *httprouter.Router

}

func (r *Routing) Init() {

	controllers.Init()
	r.Render = render.New(render.Options{Directory: "views",
		IndentJSON: true,
		Funcs: []template.FuncMap{
        {

            "str2html": func(raw string) template.HTML {
            	fmt.Println(raw)
                return template.HTML(raw)
            },
            "add": func(x,y int) int {
                return x + y
            },
            "mod": func(x,y int) int {
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

	r.Mux.GET("/"+ APIVERSION + "/application/create", ac.CreateApplication)
	r.Mux.GET("/"+ APIVERSION + "/applications/:id/delete", ac.DeleteApplication)
	r.Mux.GET("/"+ APIVERSION + "/applications/:id/update", ac.EditApplication)
	r.Mux.GET("/"+ APIVERSION + "/applications/:id", ac.ApplicationInformation)
	r.Mux.GET("/"+ APIVERSION + "/applications/:id/", ac.ApplicationInformation)

	//Container Management
	r.Mux.GET("/"+ APIVERSION + "/container/create", cc.CreateContainer)
	r.Mux.GET("/"+ APIVERSION + "/containers/:id/delete", cc.DeleteContainer)
	r.Mux.GET("/"+ APIVERSION + "/containers/:id/update", cc.EditContainer)
	r.Mux.GET("/"+ APIVERSION + "/containers/:id", cc.ContainerInformation)
	r.Mux.GET("/"+ APIVERSION + "/containers/:id/", cc.ContainerInformation)

	//Container Management
	r.Mux.POST("/"+ APIVERSION + "/node/create", nc.CreateNode)
	r.Mux.GET("/"+ APIVERSION + "/nodes/:id/delete", nc.DeleteNode)
	r.Mux.GET("/"+ APIVERSION + "/nodes/:id/update", nc.EditNode)
	r.Mux.GET("/"+ APIVERSION + "/nodes/:id", nc.NodeInformation)
	r.Mux.GET("/"+ APIVERSION + "/nodes/:id/", nc.NodeInformation)

	//User Management
	r.Mux.GET("/"+ APIVERSION + "/user/create", uc.CreateUser)
	r.Mux.GET("/"+ APIVERSION + "/users/:id/delete", uc.DeleteUser)
	r.Mux.GET("/"+ APIVERSION + "/users/:id/update", uc.EditUser)
	r.Mux.GET("/"+ APIVERSION + "/users/:id", uc.UserInformation)
	r.Mux.GET("/"+ APIVERSION + "/users/:id/", uc.UserInformation)


	r.Mux.NotFound = http.FileServer(http.Dir("public"))
}