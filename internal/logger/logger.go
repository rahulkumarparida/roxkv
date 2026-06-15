package logger

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/rahulkumarparida/roxkv/internal/utils"
)


func WriteToFile(msg string) string{
	path := utils.LogFolder()
	if path == "" {
		fmt.Printf("Error joining path")
	}

	 os.MkdirAll(path,0755)

	fileName := time.Now().Format("2006-01-02")+".txt"
	

	

	logFile := filepath.Join(path,fileName)
	file,err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // Writes appended if exists Creates if not

	_,err = file.WriteString(msg)
	
	if err != nil {
		fmt.Println("Error while reading/writing file")
		return "" 
	}
	defer file.Close()	

	return logFile

}

func ReadFromFile() (string, error){
	path := utils.LogFolder()
	if path == "" {
		fmt.Printf("No history has been recorded")
		return "" , nil
	}
	
	files, err := os.ReadDir(path)
	if err != nil {
		return "" , err
	}

	var maxTime  time.Time
	var newestFile fs.FileInfo
	for _, file := range files {
		if file.IsDir(){
			continue
		}
		info , err := file.Info()
		if err != nil {
		return "" , err
		}

		if info.ModTime().After(maxTime){
			maxTime = info.ModTime()
			newestFile = info
		}
	}

	filePath := filepath.Join(path, newestFile.Name())

	databytes , err := os.ReadFile(filePath)

	data := string(databytes)

	return data , nil

}

func SucessLog(msg string){
	messageFormat := fmt.Sprint("SUCCESS ", time.Now().Format("2006-01-02 15:04:05") ," : ",msg, " \n")
	WriteToFile(messageFormat)
}


func ErrorLog(msg string){
	messageFormat := fmt.Sprint("ERROR ",  time.Now().Format("2006-01-02 15:04:05") ," : ",msg, " \n")
	WriteToFile(messageFormat)
}


func InfoLog(msg string){
	messageFormat := fmt.Sprint("INFO ",  time.Now().Format("2006-01-02 15:04:05") ," : ",msg, " \n")
	WriteToFile(messageFormat)
}