package watcher2

import (
	"github.com/klouds/kDaemon/database"
	"github.com/klouds/kDaemon/logging"
	"github.com/klouds/kDaemon/models"
	// "errors"
	"github.com/twinj/uuid"
	"strconv"
)

var (
	TaskHandler *taskManager
)

type taskManager struct {
	node_managers map[string]NodeManager
	tasks         []Task
	jobchan       chan bool
	stopChannels  map[string]chan bool
}

func (th *taskManager) Init() {
	logging.Log("TaskHandler Init")
	if TaskHandler == nil {
		TaskHandler = &taskManager{
			node_managers: make(map[string]NodeManager),
			tasks:         make([]Task, 0),
			jobchan:       make(chan bool),
			stopChannels:  make(map[string]chan bool),
		}
	}
}

//Listens for new jobs
func (th *taskManager) Listen(stop chan bool) {
	for {
		select {
		case <-stop:
			stop <- true
			return
		case <-th.jobchan:
			//Grab first task
			task := th.tasks[0]
			//run the task
			th.Dispatch(task)

		}
	}
}

//Runs given job
func (th *taskManager) Dispatch(task Task) {
	defer th.deleteYourself(task)

	switch task.Name {

	case Launch:
		_, err := th.determineBestNodeForLaunch()
		if err != nil {
			logging.Log("There was a problem: ", err)
		}
	case Stop:

	case Delete:
	case Down:
	case Check:
	case AddNode:

	}
}

//Runs given job
func (th *taskManager) deleteYourself(task Task) {
	for index, smalltask := range th.tasks {
		if task.JobID == smalltask.JobID && len(th.tasks) > 1 {

			th.tasks = append(th.tasks[:index], th.tasks[index+1:]...)
		}
	}
}

//Adds jobs to the queue
func (th *taskManager) AddJob(name string, imageid string, containerid string) {

	newjob := Task{}
	newjob.JobID = uuid.NewV4().String()
	newjob.Name = name
	newjob.ImageID = imageid
	newjob.ContainerID = containerid

	switch name {

	case Launch:
		//Add launch command to back of queue
		th.tasks = append(th.tasks, newjob)
	case Stop:
		//Add stop command to back of queue
		th.tasks = append(th.tasks, newjob)
	case Delete:
		//Add delete command to back of queue
		th.tasks = append(th.tasks, newjob)
	case Down:
		//Add down command to front of queue
		th.tasks = append([]Task{newjob}, th.tasks...)
	case Check:
		//add check command to front of queue
		th.tasks = append([]Task{newjob}, th.tasks...)
	case AddNode:
		//add check command to front of queue
		th.tasks = append([]Task{newjob}, th.tasks...)

	}
	th.jobchan <- true
}

func (th *taskManager) determineBestNodeForLaunch() (string, error) {

	nodes, err := database.GetNodes()

	if err != nil {
		logging.Log(err)
		return "", err
	}
	idealnode := models.Node{}

	for _, value := range nodes {
		if value.State == "UP" {
			if idealnode.Id == "" {
				idealnode = value
				continue
			}
			idealcount, err := strconv.Atoi(idealnode.ContainerCount)
			if err != nil {
				logging.Log("Not a real value")
				continue
			}
			nodecount, err := strconv.Atoi(value.ContainerCount)
			if err != nil {
				logging.Log("Not a real value")
				continue
			}

			if idealcount > nodecount {
				idealnode = value
				continue
			}

		}

	}
	return "", err

}
