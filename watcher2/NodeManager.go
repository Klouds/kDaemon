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
	Node         *models.Node
	stopChannels map[string]chan bool
	jobchan      chan bool
}

//initializes the manager.
func (nm *NodeManager) Init(id string) {

	newnode, err := database.GetNode(id)
	if err != nil {
		return
	}

	maps := make(map[string]chan bool)
	nm.stopChannels = maps
	nm.Node = newnode
	nm.tasks = make([]Task, 0)
	nm.jobchan = make(chan bool)

}

//Adds jobs to the queue
func (nm *NodeManager) AddJob(task *Task) {
	stop := nm.newStopChannel(task.JobID)

	task.Stop = stop
	switch task.Name {

	case Launch:
		//Add launch command to back of queue
		nm.tasks = append(nm.tasks, *task)
	case Stop:
		//Add stop command to back of queue
		nm.tasks = append(nm.tasks, *task)
	case Delete:
		//Add delete command to back of queue
		nm.tasks = append(nm.tasks, *task)
	case Down:
		//Add down command to front of queue
		nm.tasks = append([]Task{*task}, nm.tasks...)
	case Check:
		//add check command to front of queue
		nm.tasks = append([]Task{*task}, nm.tasks...)
	case AddNode:
		//add check command to front of queue
		nm.tasks = append([]Task{*task}, nm.tasks...)

	}

	logging.Log("LENGTH : ", len(nm.tasks))
	nm.jobchan <- true
}

//Listens for new jobs
func (nm *NodeManager) Listen(stop chan bool) {
	logging.Log("I am listening ")
	for {
		select {
		case <-stop:
			stop <- true
			return
		case <-nm.jobchan:
			//Grab first task
			nm.dispatch(nm.tasks[0])
			//run the task

		}
	}
}

//Runs given job
func (nm *NodeManager) dispatch(task Task) {
	defer nm.deleteYourself(task)

	switch task.Name {

	case Launch:
		//Launch a thing

	case Stop:

	case Delete:
	case Down:
	case Check:
	case AddNode:

	}

	logging.Log("Dispatched Node job")
	//This is where the job runs
}

//deletes the job
func (nm *NodeManager) deleteYourself(task Task) {
	for index, smalltask := range nm.tasks {
		if task.JobID == smalltask.JobID && len(nm.tasks) > 1 {

			nm.tasks = append(nm.tasks[:index], nm.tasks[index+1:]...)
		}
	}
}

func (nm *NodeManager) newStopChannel(stopKey string) chan bool {
	stop := make(chan bool)
	nm.stopChannels[stopKey] = stop
	return stop
}

func (nm *NodeManager) stopForKey(key string) {
	if nm.stopChannels != nil {
		nm.stopChannels = make(map[string]chan bool)
	}

	if ch, found := nm.stopChannels[key]; found {
		ch <- true
		delete(nm.stopChannels, key)
	}
}
