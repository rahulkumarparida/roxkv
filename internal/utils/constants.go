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