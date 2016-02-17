package watcher2

import (
	"github.com/klouds/kDaemon/database"
	"github.com/klouds/kDaemon/logging"
	"github.com/klouds/kDaemon/models"
	"net"
	"time"
)

const timeout = time.Duration(5) * time.Second

func CheckNodes() ([]models.Node, error) {
	nodes, err := database.GetNodes()

	if err != nil {
		logging.Log("HC > HEALTHCHECK CANNOT START. THERE ARE NO NODES")
		return nodes, err
	}

	for index, _ := range nodes {
		//Check Node for basic ping
		conn, err := net.DialTimeout("tcp", nodes[index].DIPAddr+":"+nodes[index].DPort, timeout)
		if err != nil {
			logging.Log("HC > NODE | " + nodes[index].Name + " | IS CURRENTLY NOT ACCESSIBLE, FLAG AS DOWN")

			continue
		}

		logging.Log("HC > NODE WITH HOSTNAME | " + nodes[index].Name + " | IS HEALTHY, FLAG AS UP")

		//TODO:
		//
		//
		//
		//Send up task to node
		//
		//
		//
		conn.Close()

	}
	return nodes, nil

}

func CheckContainers() error {
	logging.Log("HC > STARTING CONTAINER CHECK")

	_, err := database.GetContainers()
	if err != nil {
		logging.Log("HC > THERE ARE NO CONTAINERS, SKIPPING HEALTHCHECK")
		return err
	}

	// for index, value := range containers {
	// 	//Send check container command
	// }

	// for index, value := range containers {
	// 	//Wait for check container response
	// }

	return err
}
