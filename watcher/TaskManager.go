package watcher

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

func (th *taskManager) CheckContainer(nodeid string, containername string) bool {
	return th.node_managers[nodeid].CheckContainer(containername)
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
	logging.Log("Listening TM")
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
				logging.Log("Job Received: ", task)
				th.Dispatch(task)

			}

		}
	}
}

//Runs given job
func (th *taskManager) Dispatch(task Task) {
	defer th.deleteYourself(&task)

	logging.Log("Dispatching job: ", task)
	if len(th.node_managers) <= 0 && task.Name != NodeUp &&
		task.Name != NodeDown && task.Name != AddNode {
		//time.Sleep(500 * time.Microsecond)
		return
	}
	switch task.Name {

	case Launch:
		node, err := th.determineBestNodeForLaunch()
		if err != nil {
			logging.Log("There was a problem: ", err)
		}

		nm := th.node_managers[node]

		if nm != nil {
			nm.AddJob(task)
		}

	case Stop:
		nm := th.node_managers[task.NodeID]
		if task.NodeID != "" && nm != nil {
			nm.AddJob(task)
		}

	//CASE DELETE WILL REMOVE ANY DATA FOR GOOD. USE AT OWN PERIL!
	case Delete:
		nm := th.node_managers[task.NodeID]
		//if container is running
		if task.NodeID != "" && nm != nil {
			stopTask := task
			stopTask.Name = Stop
			nm.AddJob(stopTask)
		}

		//DELETE ALL THE DATA
		th.deleteContainerData(task.ContainerID)

	case NodeDown:
		logging.Log("Node Down")
		nm := th.node_managers[task.NodeID]
		if task.NodeID != "" && nm != nil {
			th.node_managers[task.NodeID].AddJob(task)
		}
	case NodeUp:
		logging.Log("Node up")
		nm := th.node_managers[task.NodeID]
		if task.NodeID != "" && nm != nil {
			th.node_managers[task.NodeID].AddJob(task)
		}

	case Check:
		nm := th.node_managers[task.NodeID]
		if task.NodeID != "" && nm != nil {
			th.node_managers[task.NodeID].AddJob(task)
		}
	case AddNode:
		if task.NodeID != "" {
			//Add new node
			logging.Log("Adding node")
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

		th.nodeAddedToCluster(node.Id)
	}

	return nil
}

func (th *taskManager) nodeAddedToCluster(id string) {
	manager := NodeManager{}

	th.node_managers[id] = &manager

	manager.Init(id)
	stop := make(chan bool)
	th.stopChannels[id] = stop

	logging.Log(id)
	logging.Log(manager)
	go manager.Listen(th.stopChannels[id])
}

//deletes the job
func (th *taskManager) deleteYourself(task *Task) {

	th.tasks.Remove(task.JobID)
	task = &Task{}
	return

}

//Adds jobs to the queue
func (th *taskManager) AddJob(name string, applicationid string,
	containerid string, newname string, nodeid string) {

	newjob := &Task{}
	jobid := uuid.NewV4().String()
	newjob.JobID = jobid
	newjob.Name = name

	newjob.ApplicationID = applicationid
	newjob.ContainerID = containerid
	newjob.NewName = newname
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
