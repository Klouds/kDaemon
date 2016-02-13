package client

import (
	r "github.com/dancannon/gorethink"
	"github.com/gorilla/websocket"
	"github.com/superordinate/kDaemon/models"
)

type Handler func(*Client, interface{})

type FindHandler func(string) (Handler, bool)

type Client struct {
	Send         chan models.Message
	socket       *websocket.Conn
	findHandler  FindHandler
	Session      *r.Session
	stopChannels map[int]chan bool
	id           string
	userName     string
}

func (c *Client) NewStopChannel(stopKey int) chan bool {
	c.StopForKey(stopKey)
	stop := make(chan bool)
	c.stopChannels[stopKey] = stop
	return stop
}

func (c *Client) StopForKey(key int) {
	if ch, found := c.stopChannels[key]; found {
		ch <- true
		delete(c.stopChannels, key)
	}
}

func (client *Client) Read() {
	var message models.Message
	for {
		if err := client.socket.ReadJSON(&message); err != nil {
			break
		}
		if handler, found := client.findHandler(message.Name); found {
			handler(client, message.Data)
		}
	}
	client.socket.Close()
}

func (client *Client) Write() {
	for msg := range client.Send {
		if err := client.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	client.socket.Close()
}

func (c *Client) Close() {
	for _, ch := range c.stopChannels {
		ch <- true
	}
	close(c.Send)
}

func NewClient(socket *websocket.Conn, findHandler FindHandler,
	session *r.Session) *Client {
	return &Client{
		Send:         make(chan models.Message),
		socket:       socket,
		findHandler:  findHandler,
		Session:      session,
		stopChannels: make(map[int]chan bool),
	}
}
