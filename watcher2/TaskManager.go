package watcher2

import (
	"github.com/klouds/kDaemon/database"
	"github.com/klouds/kDaemon/logging"
	"github.com/klouds/kDaemon/models"
	cmap "github.com/streamrail/concurrent-map"
	// "errors"
	"github.com/twinj/uuid"
	"strconv"
	// "time"
)

var (
	TaskHandler *taskManager
)

type taskManager struct {
	node_managers map[string]*NodeManager
	tasks         *cmap.ConcurrentMap
	jobchan       chan bool
	stopChannels  map[string]chan bool
}

func (th *taskManager) Init() {
	logging.Log("TaskHandler Init")
	tasks := cmap.New()
	if TaskHandler == nil {
		TaskHandler = &taskManager{
			node_managers: make(map[string]*NodeManager),
			tasks:         &tasks,
			jobchan:       make(chan bool),
			stopChannels:  make(map[string]chan bool),
		}
	}

}

//Listens for new jobs
func (th *taskManager) Listen(stop chan bool) {
	//init node list
	th.initializeNodes()

	for {
		select {
		case <-stop:
			//stop <- true
			th.Shutdown()
			return
		case <-th.jobchan:
			//Grab first task
			for job := range th.tasks.IterBuffered() {

				//time.Sleep(5 * time.Microsecond)
				task := job.Val.(Task)
				th.Dispatch(task)

			}

		}
	}
}

//Runs given job
func (th *taskManager) Dispatch(task Task) {
	defer th.deleteYourself(&task)

	if len(th.node_managers) <= 0 {
		//time.Sleep(500 * time.Microsecond)
		return
	}
	switch task.Name {

	case Launch:
		node, err := th.determineBestNodeForLaunch()
		if err != nil {
			logging.Log("There was a problem: ", err)
		}

		th.node_managers[node].AddJob(task)

	case Stop:
		if task.NodeID != "" {
			th.node_managers[task.NodeID].AddJob(task)
		}

	//CASE DELETE WILL REMOVE ANY DATA FOR GOOD. USE AT OWN PERIL!
	case Delete:
		//if container is running
		if task.NodeID != "" {
			stopTask := task
			stopTask.Name = Stop
			th.node_managers[task.NodeID].AddJob(stopTask)
		}

		//DELETE ALL THE DATA
		th.deleteContainerData(task.ContainerID)

	case Down:
		if task.NodeID != "" {
			th.node_managers[task.NodeID].AddJob(task)
		}

	case Check:
		if task.NodeID != "" {
			th.node_managers[task.NodeID].AddJob(task)
		}
	case AddNode:
		if task.NodeID != "" {
			//Add new node
			th.nodeAddedToCluster(task.NodeID)
		}
	}

	return
}

func (th *taskManager) deleteContainerData(containerid string) {
	//DELETES ALL CONTAINER DATA
	//FOR HARD REMOVAL OF DATA, SHOULD BARELY EVER BE USED
	//ALL WARNINGS GO HERE.
	//
	//logging.Log("Deleting all container data for id: ", containerid)
}

func (th *taskManager) initializeNodes() error {
	//get the nodes currently in the database
	nodes, err := database.GetNodes()

	if err != nil {
		logging.Log("error getting nodes: ", err)
		return err
	}

	for _, node := range nodes {
		manager := NodeManager{}
		manager.Init(node.Id)

		th.node_managers[node.Id] = &manager
		th.nodeAddedToCluster(node.Id)
	}

	return nil
}

func (th *taskManager) nodeAddedToCluster(id string) {

	manager := th.node_managers[id]
	stop := make(chan bool)
	th.stopChannels[id] = stop

	go manager.Listen(th.stopChannels[id])
}

//deletes the job
func (th *taskManager) deleteYourself(task *Task) {

	th.tasks.Remove(task.JobID)
	task = &Task{}
	return

}

//Adds jobs to the queue
func (th *taskManager) AddJob(name string, imageid string, containerid string, nodeid string) {

	newjob := &Task{}
	jobid := uuid.NewV4().String()
	newjob.JobID = jobid
	newjob.Name = name
	newjob.ImageID = imageid
	newjob.ContainerID = containerid
	newjob.NodeID = nodeid

	th.tasks.Set(jobid, *newjob)

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
	return idealnode.Id, err

}

func (th *taskManager) Shutdown() {
	logging.Log("Shutting down")
	//Shuts down all listeners
	for key, _ := range th.node_managers {
		th.stopChannels[key] <- true
	}

	th.node_managers = make(map[string]*NodeManager)
}
