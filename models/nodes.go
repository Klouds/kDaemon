package models

import (
	"encoding/json"
	"github.com/superordinate/kDaemon/logging"
)

type Node struct {
	Id             int64  `json:"id"`
	UserID         int64  `sql:"not null;" json:"user_id"`
	Hostname       string `sql:"size:255; not null; unique;" json:"hostname"`
	DIPAddr        string `sql:"size:255; not null; unique;" json:"d_ipaddr"` //docker
	DPort          string `sql:"size:30; not null;" json:"d_port"`
	PIPAddr        string `sql:"size:30; not null; unique;" json:"p_ipaddr"` //prometheus
	PPort          string `sql:"size:255; not null;" json:"p_port"`
	ContainerCount int    `json:"container_count"`
	IsHealthy      bool   `json:"is_healthy"`
	IsEnabled      bool   `sql:"default:true" json:"is_enabled"`
	//Add location settings
	Location map[string]string `json:"location"`
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

	valid = ValidIP4(n.DIPAddr) && ValidIP4(n.PIPAddr) &&
		ValidPort(n.DPort) && ValidPort(n.PPort)

	return valid
}

func (n *Node) AddContainer() {
	n.ContainerCount = n.ContainerCount + 1
}
