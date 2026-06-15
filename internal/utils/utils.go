package utils

import (
	"log"

)


func HandleError(errmsg string, err error) bool {
	
	if err != nil {
		log.Println(errmsg, err) 
		return true         
	}
	return false 
}