package watcher

import ()

const (
	Launch   = "LAUNCH"
	Stop     = "STOP"
	Delete   = "DELETE"
	NodeDown = "NODEDOWN"
	NodeUp   = "NODEUP"
	Check    = "CHECK"
	AddNode  = "ADDNODE"
)

type Task struct {
	Dispatched    bool
	NewName       string
	JobID         string
	Name          string
	ContainerID   string
	ApplicationID string
	NodeID        string
	Stop          chan<- bool
}
