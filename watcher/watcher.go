/*					kDaemon Watcher	
	Author: 	Paul Mauviel (http://github.com/ozzadar)

	This package watches the cluster state and maintains container state across the cluster.

	Runs as a separate goroutine =)

	Responsibilities:
		- Poll for monitoring data
		- Migrate containers to ideal location
		- Launch and Destroy containers
		- Update Consul status for forwarding

*/

package watcher

import (
		"github.com/superordinate/kDaemon/models"
		"github.com/superordinate/kDaemon/database"
		"github.com/superordinate/kDaemon/logging"
		"encoding/json"
		"strings"
)

var commands = [...]string{
	"LC",  //Launch Container
	"SC",  //Shutdown Container
	"AN",  //Add Node
	"RN",  //Remove Node
	"NAC", //Not a command
}

type Job struct {
	Type 		string
	Body 		string
	InUse 		bool
	Complete 	bool	//when complete, remove job from queue
}

var queue []*Job


func MainLoop() {
	//Starts the watcher loop.
	logging.Log("Watcher started")
	for {
		RunQueue()
	}

}


//Add to queue
func AddJob(command string, object models.JSONObject) {

	for _, element := range commands {
		if element == command {
			//Valid command
			body, err := object.GetJSON()

			if err != nil {
				logging.Log(err)
				return
			}

			newjob := Job{Type: command,
								Body: body,
								InUse: false,
								Complete: false}
			queue = append(queue, &newjob)
			break
		}
	}
}


//Job Queue
func RunQueue() {
	for index, job := range queue {
		if job.Complete == true {
			job.InUse = true;
			DeleteJob(index)
			continue;
		}

		if job.Type == "LC" {
			if (job.InUse == false) {
				job.InUse = true;
				go AddContainer(job)
			}
		}
	}
}


//Commands
func AddContainer(job *Job){
	

	newcontainer := models.Container{}
	decoder := json.NewDecoder(strings.NewReader(job.Body))
	err := decoder.Decode(&newcontainer)
	if err != nil {
		logging.Log(err)
		job.Complete = true		//bad information, don't try to launch again
		job.InUse = false
		return
	}

	/* Determine node to launch on */
	id := DetermineBestNodeForLaunch()
	node, err := database.GetNode(id)
	if err != nil {
		job.Complete = true 	//bad node, so try to launch in the future
		logging.Log(err)
		job.InUse = false
		return 
	}

	/* Get the application information */
	app, err := database.GetApplication(newcontainer.ApplicationID)
	if err != nil {
		logging.Log(err)
		job.Complete = true		//Application doesn't exist, don't try to launch in the future
		job.InUse = false
		return
	}

	//compiler errors without using values right now
	node.Hostname = "rawr"
	app.Name = "fake"


	LaunchAppOnNode(app, node)

	//try to create container
		//pull if image not found
		//try to create again

	//start container


	logging.Log(newcontainer.Name)
 	job.InUse = false
	job.Complete = true
	return

}

func DeleteJob(i int) {
	queue = append(queue[:i], queue[i+1:]...)
}

func LaunchAppOnNode(app *models.Application, node *models.Node) {


}