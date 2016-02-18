package watcher2

import (
	//"github.com/klouds/kDaemon/database"
	"github.com/klouds/kDaemon/logging"
	"github.com/twinj/uuid"
)

var (
	TaskHandler *taskManager
)

type taskManager struct {
	node_managers []NodeManager
	tasks         []Task
	jobchan       chan bool
	stopChannels  map[string]chan bool
}

func (th *taskManager) Init() {
	logging.Log("TaskHandler Init")
	if TaskHandler == nil {
		TaskHandler = &taskManager{
			node_managers: make([]NodeManager, 0),
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

	logging.Log("# OF TASKS: ", len(th.tasks))

	//select node to run on

	//This is where the job runs
}

//Runs given job
func (th *taskManager) deleteYourself(task Task) {
	logging.Log("Delete yourself: ", task.JobID)
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

	if name == Launch {
		//Add launch command to back of queue
		th.tasks = append(th.tasks, newjob)
	} else if name == Stop {
		//Add stop command to back of queue
		th.tasks = append(th.tasks, newjob)
	} else if name == Delete {
		//Add delete command to back of queue
		th.tasks = append(th.tasks, newjob)
	} else if name == Down {
		//Add down command to front of queue
		th.tasks = append([]Task{newjob}, th.tasks...)
	} else if name == Check {
		//add check command to front of queue
		th.tasks = append([]Task{newjob}, th.tasks...)
	}

	th.jobchan <- true
}
