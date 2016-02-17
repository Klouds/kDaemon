package controllers

import (
	// "fmt"
	r "github.com/dancannon/gorethink"
	"github.com/klouds/kDaemon/client"
	"github.com/klouds/kDaemon/models"
	//"time"
)

const (
	NodesStop = iota
	ApplicationsStop
	ContainersStop
)

func SubscribeNodes(client *client.Client, data interface{}) {
	go func() {
		stop := client.NewStopChannel(NodesStop)
		cursor, err := r.Table("nodes").
			Changes(r.ChangesOpts{IncludeInitial: true}).
			Run(client.Session)

		if err != nil {
			client.Send <- models.Message{"error", err.Error()}
			return
		}

		changeFeedHelper(cursor, "nodes", client.Send, stop)
	}()
}

func UnsubscribeNodes(client *client.Client, data interface{}) {
	client.StopForKey(NodesStop)
}

func SubscribeApplications(client *client.Client, data interface{}) {
	go func() {
		stop := client.NewStopChannel(ApplicationsStop)
		cursor, err := r.Table("applications").
			Changes(r.ChangesOpts{IncludeInitial: true}).
			Run(client.Session)

		if err != nil {
			client.Send <- models.Message{"error", err.Error()}
			return
		}

		changeFeedHelper(cursor, "applications", client.Send, stop)
	}()
}

func UnsubscribeApplications(client *client.Client, data interface{}) {
	client.StopForKey(ApplicationsStop)
}

func SubscribeContainers(client *client.Client, data interface{}) {
	go func() {
		stop := client.NewStopChannel(ContainersStop)
		cursor, err := r.Table("containers").
			Changes(r.ChangesOpts{IncludeInitial: true}).
			Run(client.Session)
		if err != nil {
			client.Send <- models.Message{"error", err.Error()}
			return
		}
		changeFeedHelper(cursor, "containers", client.Send, stop)
	}()
}

func UnsubscribeContainers(client *client.Client, data interface{}) {
	client.StopForKey(ContainersStop)
}

func IndexPage(client *client.Client, data interface{}) {
	client.StopForKey(ContainersStop)
}

func changeFeedHelper(cursor *r.Cursor, changeEventName string,
	send chan<- models.Message, stop <-chan bool) {
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
			send <- models.Message{eventName, data}
		}
	}
}
