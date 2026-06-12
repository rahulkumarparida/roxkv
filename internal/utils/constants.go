package utils

import (
	"os"
	"path/filepath"
)


func DbFolder()  string{
	HomePath , err := os.UserHomeDir()
	dbFolder := filepath.Join(HomePath , ".roxkv" , "roxdb") 
	if HandleError("Error while fetching Home directory ", err) {
		return ""
	}
	return dbFolder

}


func LogFolder() string{
	HomePath , err := os.UserHomeDir()
	logFolder := filepath.Join(HomePath , ".roxkv" , "roxlogs") 
	if HandleError("Error while fetching Home directory ", err) {
		return ""
	}
	return logFolder
}