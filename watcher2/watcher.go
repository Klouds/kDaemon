package watcher2

import (
	//"github.com/klouds/kDaemon/database"
	"github.com/klouds/kDaemon/logging"
	"time"
)

/*
	Cluster Watcher Package

	This package is the meat and potatoes of kDaemon.
	It is designed to handle Container management and
	Health Checking to allow for 100% uptime and full
	cluster consistency.

	The watcher will be composed of several Components:

	< Task Manager >

	Handles incoming container task requests.

	TASKS:
	 	- SelectNode 		> This selects the node to run
	 						the task on.

	 	- LaunchOnNode 		> This launches the container
	 						on the given node.

	 	- StopContainer 	> This stops the container on
	 						the node it's running on.

	 	- DeleteContainer 	> This deletes the container and
	 						its data from the cluster.

	 	- FlagNodeDown		> Passes a 'DOWN' command to the
	 						appropriate Node Manager.
	 						This task is put at the beginning
	 						of the queue and should run before
	 						all other tasks.

	 	- CheckContainer 	> Passes a 'CHECK' command to the
	 						appropriate Node Manager.
	 						This task is put at the beginning
	 						of the queue and should run before
	 						all other tasks.



	< Node Manager >

	Handles tasks given to it from the Task Manager. A separate
	Node Manager is created for each Node in the system.

	TASKS:
		- LaunchContainer	> This launches the container on the
							node.
							If it fails, it returns the
							task to the Task Manager.
							On success, flags container as 'UP'
							and deletes the task.

		- StopContainer		> This stops the container on the node.
							If it fails, it returns the task to
							the Task Manager.
							On success, flags container as 'STOPPED'
							and deletes the task.

		- DeleteContainer	> This stops the container and then
							deletes it's container data from the
							Node.
							If it fails, it returns the task to
							the Task Manager.
							On success, it removes the container
							the database.

		- FlagAsDown		> This task flags the node as 'DOWN'.
							When in the 'DOWN' state, all tasks
							automatically fail and get returned
							to the Task Manager for Node Reselection.
							This task is put at the beginning
	 						of the queue and should run before
	 						all other tasks.

	 	- CheckContainer	> This task checks the health of
	 						the given container on the Node.
	 						If the container is down and not in
	 						the queue it is marked as 'DOWN' and
	 						the task returns FALSE.
	 						If the container is up, it is marked
	 						as such and the task returns TRUE.


	 < Health Checker >

	 Handles healthchecking.

	 TASKS:

	 	- CheckNodes		> Does a simple connection test to the
	 						Docker API on each Node.
	 						If a connection cannot be established, it
	 						gives the Task Manager a 'FlagAsDown'
	 						command.
	 						Keeps track of all nodes that are 'downed'

		- CheckContainers	> Sends a 'CheckContainer' command to the
							Task Manager and waits for it to complete.
							If the task fails, the container has been
							marked as 'DOWN', therefore it sends a
							'SelectNode' command, waits for it to
							return, and then sends a 'LaunchOnNode'
							command.

*/

type Watcher struct {
	HealthCheckInterval time.Duration
	lastHealthCheck     time.Time
	stopChannel         chan bool
}

func (w *Watcher) Init() {
	w.HealthCheckInterval = time.Duration(10 * time.Second)
	w.lastHealthCheck = time.Now()
	TaskHandler.Init()

	thStop := make(chan bool)
	w.stopChannel = thStop
	go TaskHandler.Listen(w.stopChannel)
	go func() {
		// count := 0

		// for {
		// 	TaskHandler.AddJob(Launch, "fake_image", "container_id", "")
		// 	count = count + 1
		// }
	}()
}

func (w *Watcher) Run(stop <-chan bool) {
	for {
		runHealthCheck := make(chan bool)

		//Go function that will run and flag a healthcheck
		//is needed every HealthCheckInterval
		go w.healthCheckTimer(runHealthCheck)

		select {
		case <-stop:
			logging.Log("ROUTINE STOPPED")

			w.stopChannel <- true
			return
		case <-runHealthCheck:
			runHealthChecks()

		}

		//check elapsed time

	}
}

func (w *Watcher) healthCheckTimer(runHealthCheck chan<- bool) {
	for {
		elapsed := time.Since(w.lastHealthCheck)

		if elapsed >= w.HealthCheckInterval {
			runHealthCheck <- true
			w.lastHealthCheck = time.Now()

		}
	}
}

func runHealthChecks() {
	logging.Log("Run HealthCheck")

	//Check Nodes
	CheckNodes()

	//Check Containers
	CheckContainers()
}
