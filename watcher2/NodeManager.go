package watcher2

import (
	"github.com/klouds/kDaemon/database"
	"github.com/klouds/kDaemon/logging"
	"github.com/klouds/kDaemon/models"
	cmap "github.com/streamrail/concurrent-map"
	//"github.com/twinj/uuid"
	// "time"
)

type NodeManager struct {
	tasks        cmap.ConcurrentMap
	Node         *models.Node
	stopChannels map[string]chan bool
	jobchan      chan bool
	dh           *dockerHandler
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
	nm.tasks = cmap.New()
	nm.jobchan = make(chan bool)

	//make a connection to the docker handler
	nm.dh, err = NewDockerHandler(newnode.DIPAddr, newnode.DPort)
	logging.Log("NEW HANDLER ON: ", newnode.DIPAddr, newnode.DPort)
	if err != nil {
		logging.Log("Failed to connect to docker endpoint")
		nm.dh = nil
	}

}

//Adds jobs to the queue
func (nm *NodeManager) AddJob(task Task) {

	nm.tasks.Set(task.JobID, task)

	nm.jobchan <- true

	return
}

//Listens for new jobs
func (nm *NodeManager) Listen(stop chan bool) {
	count := 0
	for {
		select {
		case <-stop:
			logging.Log("NM Shutting down")
			stop <- true
			return
		case <-nm.jobchan:
			iter := nm.tasks.IterBuffered()
			for job := range iter {
				//logging.Log("dispatch")
				count++
				//time.Sleep(5 * time.Microsecond)
				nm.dispatch(job.Val.(Task), count)
				//iter = nm.tasks.Iter()
			}
		}
	}
}

//Runs given job
func (nm *NodeManager) dispatch(task Task, count int) {
	defer nm.deleteYourself(&task)
	//logging.Log("LENGTH OF NM: ", nm.tasks.Count())

	switch task.Name {

	case Launch:

		nm.launchContainer(task)
	case Stop:
		// logging.Log("Stopping container ", count)
	case Down:
		// logging.Log("Downing container ", count)
	case Check:
		// logging.Log("Check container ", count)
	default:
		// logging.Log("Something else")
	}

	return
	//This is where the job runs
}

func (nm *NodeManager) launchContainer(task Task) {
	//Launch a thing
	application, err := database.GetApplication(task.ApplicationID)
	if err != nil {
		logging.Log("Application doesn't exist")
		return
	}
	container, err := database.GetContainer(task.ContainerID)

	if err != nil {
		logging.Log("Container doesn't exist")
		return
	}

	ok := nm.dh.IsImagePresent(application.DockerImage)

	if !ok {
		pullimage := nm.dh.PullImage(application.DockerImage)

		if !pullimage {
			logging.Log("cant pull image")
			return
		} else {
			logging.Log("pulling image")
			return
		}
	}

	containerexists := nm.dh.DoesContainerExist(container.Name)

	if !containerexists {
		containercreated := nm.dh.CreateContainer(container.Name, application)

		if !containercreated {
			logging.Log("Cant launch uncreated, non-existant container")
			return
		}
	}

	containerstarted := nm.dh.StartContainer(container.Name)

	if containerstarted {
		logging.Log("Successfully launched container.")
	} else {
		logging.Log("Couldn't launch the container, already started?")
	}

	container.Status = "UP"
	database.UpdateContainer(container)
}

//deletes the job
func (nm *NodeManager) deleteYourself(task *Task) {

	nm.tasks.Remove(task.JobID)
	task = &Task{}

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
