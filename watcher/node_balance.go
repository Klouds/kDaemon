package watcher

import (
	"errors"
	"github.com/klouds/kDaemon/database"
	"github.com/klouds/kDaemon/logging"
	"github.com/klouds/kDaemon/models"
)

var i = 1

//returns <0 if node doesn't exist
func DetermineBestNodeForLaunch() (*models.Node, error) {

	nodes, err := database.GetNodes()

	if err != nil {
		logging.Log(err)
		return nil, err
	}

	if len(nodes) != 0 {
		idealnode := nodes[0]

		for i := 1; i < len(nodes); i++ {
			if idealnode.IsHealthy == false && nodes[i].IsHealthy == true {
				idealnode = nodes[i]
				continue
			}

			if idealnode.ContainerCount > nodes[i].ContainerCount {
				if nodes[i].IsHealthy == true {
					idealnode = nodes[i]
				}
			}
		}
		//On launch load balancing goes here

		if idealnode.IsHealthy == false {
			return &idealnode, errors.New("ERROR")
		}

		return &idealnode, nil
	}
	return nil, nil
}
