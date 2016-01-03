package watcher

import (
	"github.com/superordinate/kDaemon/logging"
)

var i = 1

//returns <0 if node doesn't exist
func DetermineBestNodeForLaunch() int64 {
	
	if i == 1 {
		i = 3
		logging.Log("if")
	} else {
		i = 1
		logging.Log("else")
	}

	//On launch load balancing goes here

	return int64(i)
}