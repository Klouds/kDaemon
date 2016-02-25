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
	nm.dh, err = NewDockerHandler("192.168.100.25", "2375")

	if err != nil {
		logging.Log("Failed to connect to docker endpoint")
		nm.dh = nil
	}

}

//Adds jobs to the queue
func (nm *NodeManager) AddJob(task Task) {

	logging.Log("Task ::", nm)
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
		//Launch a thing
		logging.Log("Launching container ", count)
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
