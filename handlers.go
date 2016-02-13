package main

import (
	// "fmt"
	r "github.com/dancannon/gorethink"
	"github.com/superordinate/kDaemon_new/logging"
	//"time"
)

const (
	NodesStop = iota
	ApplicationsStop
	ContainersStop
)

type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

func subscribeNodes(client *Client, data interface{}) {
	logging.Log("nodes subscribed")
	go func() {
		stop := client.NewStopChannel(NodesStop)
		cursor, err := r.Table("nodes").
			Changes(r.ChangesOpts{IncludeInitial: true}).
			Run(client.session)

		if err != nil {
			client.send <- Message{"error", err.Error()}
			return
		}

		changeFeedHelper(cursor, "nodes", client.send, stop)
	}()
}

func unsubscribeNodes(client *Client, data interface{}) {
	client.StopForKey(NodesStop)
}

func subscribeApplications(client *Client, data interface{}) {
	logging.Log("applications subscribed")
	go func() {
		stop := client.NewStopChannel(ApplicationsStop)
		cursor, err := r.Table("applications").
			Changes(r.ChangesOpts{IncludeInitial: true}).
			Run(client.session)

		if err != nil {
			client.send <- Message{"error", err.Error()}
			return
		}

		changeFeedHelper(cursor, "applications", client.send, stop)
	}()
}

func unsubscribeApplications(client *Client, data interface{}) {
	client.StopForKey(ApplicationsStop)
}

func subscribeContainers(client *Client, data interface{}) {
	logging.Log("containers subscribed")
	go func() {
		stop := client.NewStopChannel(ContainersStop)
		cursor, err := r.Table("containers").
			Changes(r.ChangesOpts{IncludeInitial: true}).
			Run(client.session)
		if err != nil {
			client.send <- Message{"error", err.Error()}
			return
		}
		changeFeedHelper(cursor, "containers", client.send, stop)
	}()
}

func unsubscribeContainers(client *Client, data interface{}) {
	client.StopForKey(ContainersStop)
}

func indexPage(client *Client, data interface{}) {
	client.StopForKey(ContainersStop)
}

func changeFeedHelper(cursor *r.Cursor, changeEventName string,
	send chan<- Message, stop <-chan bool) {
	change := make(chan r.ChangeResponse)
	cursor.Listen(change)
	for {
		eventName := ""
		var data interface{}
		select {
		case <-stop:
			cursor.Close()
			return
		case val := <-change:
			if val.NewValue != nil && val.OldValue == nil {
				eventName = changeEventName + " add"
				data = val.NewValue
			} else if val.NewValue == nil && val.OldValue != nil {
				eventName = changeEventName + " remove"
				data = val.OldValue
			} else if val.NewValue != nil && val.OldValue != nil {
				eventName = changeEventName + " edit"
				data = val.NewValue
			}
			send <- Message{eventName, data}
		}
	}
}
