package watcher2

import (
	"github.com/klouds/kDaemon/database"
	"github.com/klouds/kDaemon/logging"
	"github.com/klouds/kDaemon/models"

	//"github.com/twinj/uuid"
)

//Node Task Types

type NodeManager struct {
	tasks        []Task
	node         *models.Node
	stopChannels map[string]chan bool
	jobchan      chan bool
}

//initializes the manager.
func New(id string) *NodeManager {
	newnode, err := database.GetNode(id)

	if err != nil {
		return nil
	}

	return &NodeManager{
		stopChannels: make(map[string]chan bool),
		node:         newnode,
		tasks:        make([]Task, 0),
		jobchan:      make(chan bool),
	}
}

//Adds jobs to the queue
func (nm *NodeManager) AddJob(task *Task) {

	stop := nm.newStopChannel(task.JobID)
	task.Stop = stop

	if task.Name == Launch {
		//Add launch command to back of queue
		nm.tasks = append(nm.tasks, *task)
	} else if task.Name == Stop {
		//Add stop command to back of queue
		nm.tasks = append(nm.tasks, *task)
	} else if task.Name == Delete {
		//Add delete command to back of queue
		nm.tasks = append(nm.tasks, *task)
	} else if task.Name == Down {
		//Add down command to front of queue
		nm.tasks = append([]Task{*task}, nm.tasks...)
	} else if task.Name == Check {
		//add check command to front of queue
		nm.tasks = append([]Task{*task}, nm.tasks...)
	}

	nm.jobchan <- true
}

//Listens for new jobs
func (nm *NodeManager) Listen() {
	for {
		select {
		case <-nm.jobchan:
			logging.Log("received Node job")
			nm.Dispatch()
		}
	}
}

//Runs given job
func (nm *NodeManager) Dispatch() {
	logging.Log("Dispatched Node job")
	//This is where the job runs
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
