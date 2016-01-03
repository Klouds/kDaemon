package models

import (
	"github.com/superordinate/kDaemon/common"
	"strings"
)
type Application struct {
	  Id       		int64 			`json:"id"`
	  UserID    	int64			`sql:"not null;" json:"userid"`
	  Name	 		string			`sql:"size:255; not null; unique;" json:"name"`
	  ExposedPorts	string	`json:"exposed_ports"` //docker
	  DockerImage	string			`sql:"size:255; not null;" json:"docker_image"`
	  Dependencies 	string 	`json:"dependencies"`  
	  IsEnabled		bool 			`sql:"default:true" json:"isenabled"`
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

	for _,value := range s {

		finalstring := strings.TrimSpace(value)
		valid = valid && common.ValidPort(finalstring)
	}

	return valid
}

func (n *Application) GetPorts() []string {

	s := strings.Split(n.ExposedPorts, ",")

	return s
}