package watcher2

import (
	"github.com/klouds/kDaemon/database"
	"github.com/klouds/kDaemon/logging"
	"github.com/klouds/kDaemon/models"
	cmap "github.com/streamrail/concurrent-map"
	//"github.com/twinj/uuid"
)

type NodeManager struct {
	tasks        cmap.ConcurrentMap
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
	nm.tasks = cmap.New()
	nm.jobchan = make(chan bool)

}

//Adds jobs to the queue
func (nm *NodeManager) AddJob(task Task) {

	nm.tasks.Set(task.JobID, task)

	nm.jobchan <- true

	return
}

//Listens for new jobs
func (nm *NodeManager) Listen(stop chan bool) {
	logging.Log("I am listening ")
	for {
		select {
		case <-stop:
			logging.Log("NM Shutting down")
			stop <- true
			return
		case <-nm.jobchan:
			iter := nm.tasks.Iter()
			for job := range iter {
				go nm.dispatch(job.Val.(Task))
				//iter = nm.tasks.Iter()
			}
		}
	}
}

//Runs given job
func (nm *NodeManager) dispatch(task Task) {
	defer nm.deleteYourself(&task)
	//logging.Log("LENGTH OF NM: ", nm.tasks.Count())

	switch task.Name {

	case Launch:
		//Launch a thing
		logging.Log("Launching container on: ", nm.Node.Id)
	case Stop:
		//logging.Log("Dispatched Stop job on node: ", nm.Node.Id)
	case Down:
		logging.Log("NODE IS DOWN! : ", nm.Node.Id)
	case Check:
		logging.Log("Checking Container! : ", nm.Node.Id)
	default:
		logging.Log("Something else")
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
