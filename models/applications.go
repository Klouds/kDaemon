package models

import (
	"strings"
	"encoding/json"
	"github.com/superordinate/kDaemon/logging"
)
type Application struct {
	  Id       		int64 			`json:"id"`
	  UserID    	int64			`sql:"not null;" json:"user_id"`
	  Name	 		string			`sql:"size:255; not null; unique;" json:"name"`
	  ExposedPorts	string			`json:"exposed_ports"` //docker
	  DockerImage	string			`sql:"size:255; not null;" json:"docker_image"`
	  Dependencies 	string 			`json:"dependencies"`  
	  IsEnabled		bool 			`sql:"default:true" json:"is_enabled"`
}

//Interface function
func (a *Application) GetJSON() (string, error) {
	b, err := json.Marshal(a)
    if err != nil {
        logging.Log(err)
        return "",err;
    }
    return string(b),err;
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
		valid = valid && ValidPort(value)
	}

	return valid
}

func (n *Application) GetPorts() []string {

	s := strings.Split(n.ExposedPorts, ",")

	for _,value := range s {
		value = strings.TrimSpace(value)
	}

	return s
}