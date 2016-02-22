package logging

import (
	//"bytes"
	"fmt"
	"log"
	"os"
)

func Log(v ...interface{}) {
	file, err := os.OpenFile("server.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)

	if err != nil {

	} else {
		log.SetOutput(file)
		log.Println(v...)
	}
	fmt.Println(v...)

}
