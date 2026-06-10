package persistence

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

func CreateFile(data store.Item, absoluteFilePath string)  bool{
	

	datafile , err:= os.Create(absoluteFilePath)
	if utils.HandleError("Error while creating file: ", err) {
		return  false
	}
	
	defer datafile.Close()	
	
	databytes, err := json.Marshal(data)
	if utils.HandleError("Error while Marshaling file: ", err){
		return false
	}
	
	_ , err = datafile.Write(databytes)
	if utils.HandleError("Error while writing file: ", err){
		return  false
	}

	return true
	
}



func StoreToJson(kv store.Item) bool{
	dbFolder := utils.DbFolder()
	fmt.Println("FolderPath: ", dbFolder)
	if len(dbFolder) == 0 {	
		return false
	}
	err := os.MkdirAll(dbFolder, 0755)

	if err != nil {
		utils.HandleError("Error while creating the folder", err)
		return  false
	}

	filename := kv.Key+".json"
	KeyfolderName := filepath.Join(dbFolder,filename)
			
	val := CreateFile(kv,KeyfolderName)

	if !val {
		fmt.Println("File not created")
		return	false
	}
	
	return val
	
}






func LoadJsons(ms *store.MemoryAlloc) int{
	dbFolder := utils.DbFolder()
	fmt.Println("FolderPath: ", dbFolder)
	if len(dbFolder) == 0 {	
		fmt.Println("Home directory not found")
		return 0
	}

	// Retieves all file form the path stores in an slice
	allFiles , err := os.ReadDir(dbFolder)
	
	if utils.HandleError("Error while reading the files", err) {
		return 0
	}

	// Loops through all the files to get the data 
	for _, file := range allFiles {
		if file.IsDir() {
			continue
		}

		filePath  := filepath.Join(dbFolder ,file.Name())

		databytes , err := os.ReadFile(filePath)
		if utils.HandleError("Error while reading the folder", err) {
			
			continue
		}
		var data store.Item
		err = json.Unmarshal(databytes , &data)
		store.SetKv(ms , &data)
	}

	return 1
}