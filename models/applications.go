package models

import (
	"encoding/json"
	"github.com/klouds/kDaemon/logging"
	"strings"
)

type Application struct {
	Id           string `json:"id,omitempty" gorethink:"id,omitempty"`
	UserID       string `json:"user_id" gorethink:"user_id"`
	Name         string `json:"name" gorethink:"name"`
	ExposedPorts string `json:"exposed_ports" gorethink:"exposed_ports"` //docker
	DockerImage  string `json:"docker_image" gorethink:"docker_image"`
	Dependencies string `json:"dependencies" gorethink:"dependencies"`
	IsEnabled    bool   `json:"is_enabled" gorethink:"is_enabled"`
}

//Interface function
func (a *Application) GetJSON() (string, error) {
	b, err := json.Marshal(a)
	if err != nil {
		logging.Log(err)
		return "", err
	}
	return string(b), err
}

func (n *Application) AddPort(text string) {
	if len(n.ExposedPorts) < 1 {
		n.ExposedPorts = text
		return
	}
	n.ExposedPorts = n.ExposedPorts + "," + text
}

func (n *Application) AddDependency(text string) {
	if len(n.Dependencies) < 1 {
		n.Dependencies = text
		return
	}
	n.Dependencies = n.Dependencies + "," + text
}

func (n *Application) Validate() bool {
	valid := true

	s := n.GetPorts()

	for _, value := range s {
		valid = valid && ValidPort(value)
	}

	return valid
}

func (n *Application) GetPorts() []string {

	s := strings.Split(n.ExposedPorts, ",")

	for _, value := range s {
		value = strings.TrimSpace(value)
	}

	return s
}

func (n *Application) MergeChanges(newApp *Application) *Application {

	newapp := Application{}

	newapp = *n
	newapp.Id = n.Id

	if n.UserID != newApp.UserID && newApp.UserID != "" {
		newapp.UserID = newApp.UserID
	}

	if n.Name != newApp.Name && newApp.Name != "" {
		newapp.Name = newApp.Name
	}
	if n.ExposedPorts != newApp.ExposedPorts && newApp.ExposedPorts != "" {
		newapp.ExposedPorts = newApp.ExposedPorts
	}
	if n.DockerImage != newApp.DockerImage && newApp.DockerImage != "" {
		newapp.DockerImage = newApp.DockerImage
	}

	if n.Dependencies != newApp.Dependencies && newApp.Dependencies != "" {
		newapp.Dependencies = newApp.Dependencies
	}

	return &newapp
}
