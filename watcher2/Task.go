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
	JobID       string
	Name        string
	ContainerID string
	ImageID     string
	Stop        chan<- bool
}
