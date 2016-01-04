package watcher

import (
	"github.com/superordinate/kDaemon/logging"
	"github.com/superordinate/kDaemon/database"
	"github.com/superordinate/kDaemon/models"
)

var i = 1

//returns <0 if node doesn't exist
func DetermineBestNodeForLaunch() (*models.Node, error) {
	
	nodes, err := database.GetNodes()
	logging.Log("ALLNODES")
	logging.Log(nodes)

	if err != nil {
		logging.Log(err)
		return nil,err
	}

	idealnode := nodes[0]

	logging.Log(idealnode)

	for i:= 1; i<len(nodes); i++ {
		if idealnode.ContainerCount > nodes[i].ContainerCount {

			logging.Log("COMPARING AGAINST")
			logging.Log(nodes[i])
			idealnode = nodes[i]
		}
	}

	logging.Log("IDEAL NODE")
	logging.Log(idealnode)
	//On launch load balancing goes here

	return &idealnode, nil
}