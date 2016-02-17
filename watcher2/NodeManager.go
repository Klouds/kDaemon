package watcher2

import (
	"github.com/klouds/kDaemon/database"
	"github.com/klouds/kDaemon/logging"
	"github.com/klouds/kDaemon/models"
	"github.com/twinj/uuid"
)

//Node Task Types
const (
	Launch = "LAUNCH"
	Stop   = "STOP"
	Delete = "DELETE"
	Down   = "DOWN"
	Check  = "CHECK"
)

type NodeJob struct {
	JobID       string
	Name        string
	ContainerID string
	ImageID     string
	Stop        chan<- bool
}

type NodeManager struct {
	jobs         []NodeJob
	node         *models.Node
	stopChannels map[string]chan bool
}

//initializes the manager.
func NewNodeManager(id string) *NodeManager {
	newnode, err := database.GetNode(id)

	if err != nil {
		return nil
	}

	return &NodeManager{
		stopChannels: make(map[string]chan bool),
		node:         newnode,
		jobs:         make([]NodeJob, 0),
	}
}

//Adds jobs to the queue
func (nm *NodeManager) AddJob(name string, containerid string,
	imageid string, stop chan<- bool) {

	newjob := NodeJob{}
	newjob.JobID = uuid.NewV4().String()
	newjob.Name = name
	newjob.ImageID = imageid
	newjob.ContainerID = containerid
	nm.newStopChannel(newjob.JobID)
	newjob.Stop = stop

	if name == Launch {
		//Add launch command to back of queue
		nm.jobs = append(nm.jobs, newjob)
	} else if name == Stop {
		//Add stop command to back of queue
		nm.jobs = append(nm.jobs, newjob)
	} else if name == Delete {
		//Add delete command to back of queue
		nm.jobs = append(nm.jobs, newjob)
	} else if name == Down {
		//Add down command to front of queue
		nm.jobs = append([]NodeJob{newjob}, nm.jobs...)
	} else if name == Check {
		//add check command to front of queue
		nm.jobs = append([]NodeJob{newjob}, nm.jobs...)
	}

	logging.Log(nm.jobs)
}

func (nm *NodeManager) newStopChannel(stopKey string) chan bool {
	nm.stopForKey(stopKey)
	stop := make(chan bool)
	nm.stopChannels[stopKey] = stop
	return stop
}

func (nm *NodeManager) stopForKey(key string) {
	if ch, found := nm.stopChannels[key]; found {
		ch <- true
		delete(nm.stopChannels, key)
	}
}
