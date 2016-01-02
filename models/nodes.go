package models

import (

)

type Node struct {
	  Id       	int64 	`json:"id"`
	  UserID    int64	`sql:"not null;" json:"userid"`
	  Hostname 	string	`sql:"size:255; not null; unique;" json:"hostname"`
	  DIPAddr	string 	`sql:"size:255; not null; unique;" json:"dipaddr"` //docker
	  DPort		string	`sql:"size:30; not null;" json:"dport"`
	  PIPAddr 	string 	`sql:"size:30; not null; unique;" json:"pipaddr"`  //prometheus
	  PPort 	string 	`sql:"size:255; not null;" json:"pport"`
	  IsEnabled	bool 	`sql:"default:true" json:"isenabled"`
}