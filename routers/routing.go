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

	r.Mux.POST("/"+ APIVERSION + "/application/create", ac.CreateApplication)
	r.Mux.DELETE("/"+ APIVERSION + "/applications/:id/delete", ac.DeleteApplication)
	r.Mux.PUT("/"+ APIVERSION + "/applications/:id/update", ac.EditApplication)
	r.Mux.GET("/"+ APIVERSION + "/applications/:id", ac.ApplicationInformation)
	r.Mux.GET("/"+ APIVERSION + "/applications/:id/", ac.ApplicationInformation)

	//Container Management
	r.Mux.POST("/"+ APIVERSION + "/container/create", cc.CreateContainer)
	r.Mux.DELETE("/"+ APIVERSION + "/containers/:id/delete", cc.DeleteContainer)
	r.Mux.PUT("/"+ APIVERSION + "/containers/:id/update", cc.EditContainer)
	r.Mux.GET("/"+ APIVERSION + "/containers/:id", cc.ContainerInformation)
	r.Mux.GET("/"+ APIVERSION + "/containers/:id/", cc.ContainerInformation)

	//Container Management
	r.Mux.POST("/"+ APIVERSION + "/node/create", nc.CreateNode)
	r.Mux.DELETE("/"+ APIVERSION + "/nodes/:id/delete", nc.DeleteNode)
	r.Mux.PUT("/"+ APIVERSION + "/nodes/:id/update", nc.EditNode)
	r.Mux.GET("/"+ APIVERSION + "/nodes/:id", nc.NodeInformation)
	r.Mux.GET("/"+ APIVERSION + "/nodes/:id/", nc.NodeInformation)

	//User Management
	r.Mux.POST("/"+ APIVERSION + "/user/create", uc.CreateUser)
	r.Mux.DELETE("/"+ APIVERSION + "/users/:id/delete", uc.DeleteUser)
	r.Mux.PUT("/"+ APIVERSION + "/users/:id/update", uc.EditUser)
	r.Mux.GET("/"+ APIVERSION + "/users/:id", uc.UserInformation)
	r.Mux.GET("/"+ APIVERSION + "/users/:id/", uc.UserInformation)


	r.Mux.NotFound = http.FileServer(http.Dir("public"))
}