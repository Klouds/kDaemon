package common

import (

    "regexp"
    "strings"
    "strconv"

 )
//Validation functions

func ValidIP4(ipAddr string) bool {
	ipAddr = strings.Trim(ipAddr, " ")

	 re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
     
     if re.MatchString(ipAddr) {
             return true
     }
     return false
}

func ValidPort(port string) bool {
	port = strings.Trim(port, " ")

	i, err := strconv.Atoi(port)

	if err != nil {
		return false
	}

	if i < 0 || i > 65535 {
		return false
	} 

	return true

}