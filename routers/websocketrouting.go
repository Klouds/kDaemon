package routers

import (
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/gorilla/websocket"
	"github.com/superordinate/kDaemon/client"
	"github.com/superordinate/kDaemon/controllers"
	"github.com/superordinate/kDaemon/database"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WebSocketRouter struct {
	Router *wsrouter
}

type wsrouter struct {
	rules   map[string]client.Handler
	session *r.Session
}

func (ws_r *WebSocketRouter) Init() {

	//creates the new router
	ws_r.Router = &wsrouter{
		rules:   make(map[string]client.Handler),
		session: database.Session,
	}

	ws_r.Router.Handle("nodes subscribe", controllers.SubscribeNodes)
	ws_r.Router.Handle("nodes unsubscribe", controllers.UnsubscribeNodes)

	ws_r.Router.Handle("applications subscribe", controllers.SubscribeApplications)
	ws_r.Router.Handle("applications unsubscribe", controllers.UnsubscribeApplications)

	ws_r.Router.Handle("containers subscribe", controllers.SubscribeContainers)
	ws_r.Router.Handle("containers unsubscribe", controllers.UnsubscribeContainers)
}

func (r *wsrouter) Handle(msgName string, handler client.Handler) {
	r.rules[msgName] = handler
}

func (r *wsrouter) FindHandler(msgName string) (client.Handler, bool) {
	handler, found := r.rules[msgName]
	return handler, found
}

func (e *wsrouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	client := client.NewClient(socket, e.FindHandler, e.session)
	defer client.Close()
	go client.Write()
	client.Read()

}
