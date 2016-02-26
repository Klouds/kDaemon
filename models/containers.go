package models

import (
	"encoding/json"
	"github.com/klouds/kDaemon/logging"
)

type Container struct {
	Id            string  `json:"id,omitempty" gorethink:"id,omitempty"`
	NodeID        string  `json:"node_id" gorethink:"node_id"`
	ApplicationID string  `json:"application_id" gorethink:"application_id"`
	UserID        string  `json:"user_id" gorethink:"user_id"`
	Name          string  `json:"name" gorethink:"name"`
	ContainerID   string  `json:"container_id" gorethink:"container_id"`
	Balance       float64 `json:"balance" gorethink:"balance"`
	Status        string  `json:"status" gorethink:"status"`
	AccessLink    string  `json:"access_link" gorethink:"access_link"`
}

//Interface function
func (c Container) GetJSON() (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		logging.Log(err)
		return "", err
	}
	return string(b), err
}

func (c *Container) MergeChanges(container *Container) *Container {

	newcontainer := Container{}

	newcontainer = *c
	newcontainer.Id = c.Id

	if c.ApplicationID != container.ApplicationID && container.ApplicationID != "" {
		newcontainer.ApplicationID = container.ApplicationID
	}

	if c.Name != container.Name && container.Name != "" {
		newcontainer.Name = container.Name
	}

	if c.Status != container.Status && container.Status != "" {
		newcontainer.Status = container.Status
	}

	return &newcontainer
}
