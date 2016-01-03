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
		"github.com/superordinate/kDaemon/logging"
		"strconv"
)

/* Job Commands. For queueing up actions on the cluster.*/
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
				job.InUse = true
				go AddContainer(job)
			}
		}
	}
}



func DeleteJob(i int) {
	index := strconv.Itoa(i)
	logging.Log("Deleting job: " + queue[i].Type+ " at index " + index)
	queue = append(queue[:i], queue[i+1:]...)
}



