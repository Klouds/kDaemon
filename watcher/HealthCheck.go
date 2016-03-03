package watcher

import (
	"errors"
	"github.com/klouds/kDaemon/database"
	"github.com/klouds/kDaemon/logging"
	"github.com/klouds/kDaemon/models"
	"net"
	"strconv"
	"time"
)

const timeout = time.Duration(5) * time.Second

func CheckNodes() ([]models.Node, error) {
	nodes, err := database.GetNodes()

	if err != nil {
		return nodes, err
	}

	for _, node := range nodes {

		//Check Node for basic ping
		conn, err := net.DialTimeout("tcp", node.DIPAddr+":"+node.DPort, timeout)
		if err != nil {
			TaskHandler.AddJob(NodeDown, "", "", "", node.Id)
			continue
		}
		TaskHandler.AddJob(NodeUp, "", "", "", node.Id)

		conn.Close()

	}

	//As of this moment, the watcher will flag nodes as up or down.
	//I believe this is all we need for the actual node part. On to containers!
	return nodes, nil
}

/*
	The general process of how we want to deal with container health checking:
	1) If node is down, migrate all containers to an available node
		For now, this will loop through and mark affected containers as down

	2) We then want to loop through the containers on each healthy node and
		marking any unaccessible containers as down.

	3) We then loop through the containers marked as Down and relaunch them.

	***** PROGRAMMER'S NOTE *****
	*
	* This may require an addition of a "Should be up" type variable
	* which will be used to ignore containers that have been created
	* but not yet marked for launch.
	*
	* This will give us the added benefit of being able to keep containers
	* in our system though disabling them (For example, if payment
	* isn't successful)
	*
	* ***************************

*/
func CheckContainers() error {

	//Here, we get all the nodes
	nodes, err := database.GetNodes()

	if err != nil {
		return errors.New("No Nodes.")
	}
	//We should probably add a function to our DB that will return
	//all containers on a specific node.
	//
	//this loop will relaunch any containers that are down
	for _, node := range nodes {
		containers, err := database.GetContainersOnNode(node.Id)

		if err != nil {
			logging.Log("Node has no containers")
			continue
		}
		//for each container on the node
		for _, container := range containers {
			//if the node is down, mark container as down
			if node.State == "DOWN" {
				container.Status = "DOWN"
				container.AccessLink = ""
			} else {
				//if the node is up, check if the container is accessible.
				isup := TaskHandler.CheckContainer(node.Id, container.Name)

				if !isup {
					container.Status = "DOWN"
					container.AccessLink = ""
					TaskHandler.AddJob(Launch,
						container.ApplicationID,
						container.Id,
						container.Name,
						"")
				} else {
					container.Status = "UP"
				}
			}

			database.UpdateContainer(&container)
		}

	}

	return nil
}

//This will rebalance containers based on whether or not
//there is an imbalance in load.
func Rebalance() error {
	//Lets grab all our nodes
	nodes, err := database.GetNodesByState("UP")

	if err != nil || len(nodes) == 0 {
		return errors.New("No Nodes.")
	}

	//Let's also grab all our containers
	containers, err := database.GetContainers()

	if err != nil {
		logging.Log("No containers")
		return errors.New("No containers")
	}

	//we're doing a simple round robin balancing at the moment
	//so we'll take our container-per-node and average out the load

	if len(nodes) <= 0 {
		return errors.New("No nodes up")
	}
	containerspernode := int(len(containers) / len(nodes))
	//This will ensure that rounding errors dont fuck things up.
	containerscounted := 0
	containerstobemoved := []models.Container{}

	//we'll then go through each node and remove any containers and store
	//them for relaunch

	for _, node := range nodes {
		//for each container on the node
		containers, err := database.GetContainersOnNode(node.Id)
		if err != nil {
			logging.Log("Node has no containers")
			continue
		}
		for i, container := range containers {
			//Count this container
			containerscounted = containerscounted + 1

			if i > containerspernode {
				//if the current container is past the expect count
				//
				//remove the container from the node
				//Then add the container to the list of containers to
				//be relaunched
				containerstobemoved = append(containerstobemoved, container)
			}

		}
	}

	//Once we've collected containers to be moved
	//let's move them
	//
	for _, container := range containerstobemoved {
		TaskHandler.AddJob(Stop,
			"",
			container.Id,
			container.Name,
			container.NodeID)
		TaskHandler.AddJob(Launch,
			container.ApplicationID,
			container.Id,
			container.Name,
			"")
	}

	return nil
}

func RecountContainers() {
	containermap := make(map[string]int)
	//Let's also grab all our containers
	containers, err := database.GetContainers()

	if err != nil {
		logging.Log("No containers")
		return
	}

	for _, container := range containers {
		containermap[container.NodeID] = containermap[container.NodeID] + 1
	}

	if len(containermap) == 0 {
		return
	}
	for key, value := range containermap {
		node, err := database.GetNode(key)
		if err != nil {
			logging.Log("key doesnt exist")
			continue
		}

		node.ContainerCount = strconv.Itoa(value)

		database.UpdateNode(node)

	}
}

/*
	After we have preformed our "Services" health check, we're then
	going to perform a few node cleanup routines designed to optimize
	disk space.
*/
