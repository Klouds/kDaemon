package models

import (
	"encoding/json"
	"github.com/superordinate/kDaemon/logging"
)

type Node struct {
	Id             string `json:"id,omitempty" gorethink:"id,omitempty"`
	UserID         string `json:"user_id" gorethink:"userid"`
	Name           string `json:"name" gorethink:"name"`
	DIPAddr        string `json:"d_ipaddr" gorethink:"d_ipaddr"` //docker
	DPort          string `json:"d_port" gorethink:"d_port"`
	ContainerCount int    `json:"container_count" gorethink:"container_count"`
	IsHealthy      bool   `json:"is_healthy" gorethink:"is_healthy"`
	IsEnabled      bool   `json:"is_enabled" gorethink:"is_enabled"`
}

//Interface function
func (n *Node) GetJSON() (string, error) {
	b, err := json.Marshal(n)
	if err != nil {
		logging.Log(err)
		return "", err
	}
	return string(b), err
}

func (n *Node) Validate() bool {
	valid := true

	valid = ValidIP4(n.DIPAddr) &&
		ValidPort(n.DPort)

	return valid
}

func (n *Node) AddContainer() {
	n.ContainerCount = n.ContainerCount + 1
}
