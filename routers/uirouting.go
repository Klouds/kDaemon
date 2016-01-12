package routers

import (
	"net/http"
	"html/template"
	"gopkg.in/unrolled/render.v1"
	"github.com/julienschmidt/httprouter"
	"github.com/superordinate/kDaemon/controllers"
)

type UIRouting struct {

	Render *render.Render
	Mux *httprouter.Router
}


func (r *UIRouting) Init() {

	r.Render = render.New(render.Options{Directory: "views",
		IndentJSON: true,
		Funcs: []template.FuncMap{
        {

            "str2html": func(raw string) template.HTML {
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


	uic := &controllers.UIController{Render: r.Render}

	r.Mux.GET("/", uic.Index)

	r.Mux.NotFound = http.FileServer(http.Dir("public"))
}