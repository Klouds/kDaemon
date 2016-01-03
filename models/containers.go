package models

import (
	"encoding/json"
	"github.com/superordinate/kDaemon/logging"
)

type Container struct {
	  Id       		int64 	`json:"id"`
	  NodeID		int64	`json:"node_id"`
	  ApplicationID int64 	`json:"application_id"`
	  UserID 		int64	`json:"user_id"`
	  Name	 		string	`sql:"size:255; not null; unique;" json:"name"`
	  ContainerID	string 	`sql:"size:255; not null; unique;" json:"container_id"` 
	  Balance		float64	`json:"balance"`
	  Status 		string	`json:"status"`
	  IsEnabled		bool 	`sql:"default:true" json:"is_enabled"`
}

//Interface function
func (c Container) GetJSON() (string, error) {
	b, err := json.Marshal(c)
    if err != nil {
        logging.Log(err)
        return "",err;
    }
    return string(b),err;
}