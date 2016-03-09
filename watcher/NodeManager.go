package watcher

import (
	docker "github.com/fsouza/go-dockerclient"
	"github.com/klouds/kDaemon/database"
	"github.com/klouds/kDaemon/logging"
	"github.com/klouds/kDaemon/models"
	cmap "github.com/streamrail/concurrent-map"
	//"github.com/twinj/uuid"
	// "time"
	"strconv"
)

type NodeManager struct {
	tasks        cmap.ConcurrentMap
	Node         *models.Node
	stopChannels map[string]chan bool
	jobchan      chan bool
	dh           *dockerHandler
	isDown       bool
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
	nm.isDown = false

	nm.connectToDocker()
}

//make a connection to the docker handler
func (nm *NodeManager) connectToDocker() {
	dh, err := NewDockerHandler(nm.Node.DIPAddr, nm.Node.DPort)

	nm.dh = dh
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
				go nm.dispatch(job.Val.(Task), count)
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

		if !nm.isDown {
			nm.launchContainer(task)
		}
	case Stop:
		//Stop job...
		//
		// I think this a good point to add a stop container task to
		// our docker handler
		if !nm.isDown {
			nm.stopContainer(task)
		}
	case NodeDown:
		//This will flag the node as down, all subsequent tasks
		//fail and are returned to the task manager

		node, err := database.GetNode(nm.Node.Id)
		if err != nil {
			logging.Log("there was an error")
			return
		}
		nm.Node = node
		nm.isDown = true
		node.State = "DOWN"
		node.ContainerCount = "0"
		database.UpdateNode(node)
		nm.requeuejobs()
	case NodeUp:
		nm.isDown = false

		nm.connectToDocker()

		node, err := database.GetNode(nm.Node.Id)
		if err != nil {
			logging.Log("there was an error")
			return
		}
		nm.Node = node

		node.State = "UP"
		database.UpdateNode(node)

	case Check:
		// logging.Log("Check container ", count)
	default:
		// logging.Log("Something else")
	}

	return
	//This is where the job runs
}

//the node has been flagged as down and jobs will no longer complete.
//reschedule and delete please
func (nm *NodeManager) requeuejobs() {
	defer nm.deleteAllTasks()

	iter := nm.tasks.IterBuffered()
	for job := range iter {
		task := job.Val.(Task)
		switch task.Name {
		case Launch:
			TaskHandler.AddJob(task.Name,
				task.ApplicationID,
				task.ContainerID,
				task.NewName,
				"")
			continue
		default:

		}

	}
}

func (nm *NodeManager) deleteAllTasks() {
	//delete all the tasks in the queue
	iter := nm.tasks.IterBuffered()
	for job := range iter {
		task := job.Val.(Task)
		nm.tasks.Remove(task.JobID)
		task = Task{}
	}
}

func (nm *NodeManager) stopContainer(task Task) {
	ok := nm.dh.StopContainer(task.NewName)

	if ok {

		origcount, err := strconv.Atoi(nm.Node.ContainerCount)

		if err != nil {
			logging.Log("Not a valid container count")
			return
		}

		nm.Node.ContainerCount = strconv.Itoa(origcount - 1)
		database.UpdateNode(nm.Node)
	}

	container, err := database.GetContainer(task.ContainerID)

	if err != nil {
		logging.Log("CONTAINER DOESNT EXIST")
	}
	container.Status = "DOWN"

	database.UpdateContainer(container)
	return
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
		containercreated := nm.dh.CreateContainer(container.Name,
			application, container.GetEnvironmentVariables())

		if !containercreated {
			logging.Log("Cant launch uncreated, non-existant container")
			return
		}
	}

	containerstarted := nm.dh.StartContainer(container.Name)

	if containerstarted {

		origcount, err := strconv.Atoi(nm.Node.ContainerCount)

		if err != nil {
			logging.Log("Not a valid container count")
		}

		nm.Node.ContainerCount = strconv.Itoa(origcount + 1)
		database.UpdateNode(nm.Node)
	} else {

	}

	//inspect the container
	cont := nm.dh.InspectContainer(container.Name)

	if cont == nil {
		logging.Log("Container inspection failed. Node must be down")
	}

	ports := cont.NetworkSettings.Ports
	portbindings := []docker.PortBinding{}

	for _, value := range ports {
		portbindings = append(value, portbindings...)
	}

	if len(portbindings) > 0 {
		container.AccessLink = "http://" + nm.Node.DIPAddr + ":" + portbindings[0].HostPort
	}
	container.Status = "UP"
	container.NodeID = nm.Node.Id

	database.UpdateContainer(container)
}

//deletes the job
func (nm *NodeManager) deleteYourself(task *Task) {

	if !nm.isDown {
		nm.tasks.Remove(task.JobID)
		task = &Task{}
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

func (nm *NodeManager) CheckContainer(containername string) bool {
	return nm.dh.CheckContainerIsRunning(containername)
}
