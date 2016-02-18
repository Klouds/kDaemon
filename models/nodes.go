package models

import (
	"encoding/json"
	"github.com/klouds/kDaemon/logging"
)

type Node struct {
	Id             string `json:"id,omitempty" gorethink:"id,omitempty"`
	UserID         string `json:"user_id" gorethink:"userid"`
	Name           string `json:"name" gorethink:"name"`
	DIPAddr        string `json:"d_ipaddr" gorethink:"d_ipaddr"` //docker
	DPort          string `json:"d_port" gorethink:"d_port"`
	ContainerCount string `json:"container_count" gorethink:"container_count"`
	State          string `json:"state" gorethink:"state"`
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

func (n *Node) MergeChanges(newNode *Node) *Node {

	newnode := Node{}

	newnode = *n
	newnode.Id = n.Id

	if n.UserID != newNode.UserID && newNode.UserID != "" {
		newnode.UserID = newNode.UserID
	}

	if n.Name != newNode.Name && newNode.Name != "" {
		newnode.Name = newNode.Name
	}
	if n.DIPAddr != newNode.DIPAddr && newNode.DIPAddr != "" {
		newnode.DIPAddr = newNode.DIPAddr
	}
	if n.DPort != newNode.DPort && newNode.DPort != "" {
		newnode.DPort = newNode.DPort
	}

	if n.DPort != newNode.DPort && newNode.DPort != "" {
		newnode.DPort = newNode.DPort
	}

	if n.State != newNode.State && newNode.State != "" {
		newnode.State = newNode.State
	}

	if n.ContainerCount != newNode.ContainerCount && newNode.ContainerCount != "" {
		newnode.ContainerCount = newNode.ContainerCount
	}
	return &newnode
}
