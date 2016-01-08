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
		"time"
)

const HC_INTERVAL = time.Duration(15) * time.Second

/* Job Commands. For queueing up actions on the cluster.*/
var commands = [...]string{
	
	"SC",  //Shutdown Container
	"HC",  //Performs a global health check
	"RC",  //Removes a container from a node
	"LC",  //Launch Container on a node

	"NAC", //Not a command
}

type Job struct {
	Type 		string
	Body 		string
	InUse 		bool
	Complete 	bool	//when complete, remove job from queue
}

//The job queue
var queue []*Job


func MainLoop() {

	//Starts the watcher loop.
	logging.Log("Watcher started")
	go ScheduleHealthCheck(HC_INTERVAL)

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
	currentTime := time.Now().Local()
	for index, job := range queue {
		if job.Complete == true {
			job.InUse = true;
			DeleteJob(index)
			continue;
		}

		if job.Type == "LC" {
			if (job.InUse == false) {
				job.InUse = true
				logging.Log("LC > LAUNCHING CONTAINER ON NODE: " + currentTime.String())
				go AddContainer(job)
			}
		}

		if job.Type == "HC" {
			if (job.InUse == false) {
				
				job.InUse = true
				logging.Log("HC > PERFORMING HEALTH CHECK AT " + currentTime.String())
				go PerformHealthCheck(job)
			}
		}

		if job.Type == "RC" {
			if (job.InUse == false) {
				currentTime := time.Now().Local()
				job.InUse = true
				logging.Log("RC > REMOVING CONTAINER AT " + currentTime.String())
				go RemoveContainer(job)
			}
		}
	}
}



func DeleteJob(i int) {
	index := strconv.Itoa(i)
	logging.Log("Deleting job: " + queue[i].Type+ " at index " + index)


	queue = append(queue[:i], queue[i+1:]...)
}

func ScheduleHealthCheck(interval time.Duration) {
	for _ = range time.Tick(interval) {
		//Tick
		
		newjob := &Job {
			Type: "HC",
			Body: "{}",
			InUse: false,
			Complete: false,
		}

		queue = append(queue, newjob)

	}
}

