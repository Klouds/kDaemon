package watcher2

import ()

const (
	Launch  = "LAUNCH"
	Stop    = "STOP"
	Delete  = "DELETE"
	Down    = "DOWN"
	Check   = "CHECK"
	AddNode = "ADDNODE"
)

type Task struct {
	Dispatched  bool
	JobID       string
	Name        string
	ContainerID string
	ImageID     string
	NodeID      string
	Stop        chan<- bool
}
