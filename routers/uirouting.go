package routers

import (
	"github.com/julienschmidt/httprouter"
	"github.com/superordinate/kDaemon/controllers"
	"gopkg.in/unrolled/render.v1"
	"html/template"
	"net/http"
)

type UIRouting struct {
	Render *render.Render
	Mux    *httprouter.Router
}

func (r *UIRouting) Init() {

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

	uic := &controllers.UIController{Render: r.Render}

	//Testing auto login with git push here
	r.Mux.GET("/", uic.Index)

	//Node routing
	r.Mux.GET("/nodes", uic.NodeIndex)
	r.Mux.GET("/nodes/create", uic.CreateNode)

	//ApplicationRouting
	r.Mux.GET("/applications", uic.AppIndex)

	//Container Routing
	r.Mux.GET("/containers", uic.ContainerIndex)

	r.Mux.NotFound = http.FileServer(http.Dir("public"))
}
